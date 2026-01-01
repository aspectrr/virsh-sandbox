#!/usr/bin/env python3
"""Post-process generated SDK for better quality and add unified client with flattened parameters."""

import re
from pathlib import Path
from dataclasses import dataclass
from typing import Optional


@dataclass
class FieldInfo:
    name: str
    type_hint: str
    default: str
    description: str


@dataclass
class MethodInfo:
    name: str
    path_params: list  # List of (name, type, default)
    request_type: Optional[str]
    return_type: str
    docstring: str


def parse_model_fields(model_path: Path) -> list[FieldInfo]:
    """Parse a Pydantic model file to extract field information."""
    content = model_path.read_text()
    fields = []

    # Split content into lines for easier processing
    lines = content.split('\n')

    # Track if we're inside the class definition
    in_class = False
    class_indent = 0

    i = 0
    while i < len(lines):
        line = lines[i]
        stripped = line.strip()

        # Detect class definition
        if re.match(r'^class \w+\(BaseModel\):', stripped):
            in_class = True
            # Find the indentation level of class body
            class_indent = len(line) - len(line.lstrip()) + 4  # class body is indented
            i += 1
            continue

        # If we're in a class, look for field definitions
        if in_class:
            current_indent = len(line) - len(line.lstrip()) if line.strip() else 0

            # Check if we've left the class (less indentation or new class/def at same level)
            if stripped and current_indent < class_indent and not stripped.startswith('#'):
                if stripped.startswith('class ') or stripped.startswith('def '):
                    in_class = False
                    i += 1
                    continue

            # Stop parsing fields when we hit a method definition inside the class
            # Methods mark the end of field definitions in Pydantic models
            if stripped.startswith('def ') or stripped.startswith('async def '):
                in_class = False
                i += 1
                continue

            # Also stop if we hit model_config (usually comes after fields)
            if stripped.startswith('model_config'):
                in_class = False
                i += 1
                continue

            # Skip empty lines, comments, docstrings, and special attributes
            if (not stripped or
                stripped.startswith('#') or
                stripped.startswith('"""') or
                stripped.startswith("'''") or
                stripped.startswith('__')):
                i += 1
                continue

            # Match field definition with Field(): name: Type = Field(...)
            field_match = re.match(r'^(\w+):\s*(.+?)\s*=\s*Field\(', stripped)
            if field_match and not field_match.group(1).startswith('_'):
                field_name = field_match.group(1)
                type_hint = field_match.group(2).strip()

                # Extract inner type from Optional[]
                optional_match = re.match(r'Optional\[(.+)\]', type_hint)
                if optional_match:
                    type_hint = optional_match.group(1)

                # Find the description in this or following lines
                description = ""
                full_field_text = stripped

                # Collect the full Field() definition which may span multiple lines
                paren_count = stripped.count('(') - stripped.count(')')
                j = i + 1
                while paren_count > 0 and j < len(lines):
                    full_field_text += lines[j]
                    paren_count += lines[j].count('(') - lines[j].count(')')
                    j += 1

                # Extract description from the full field text
                desc_match = re.search(r'description\s*=\s*["\'](.+?)["\']', full_field_text)
                if desc_match:
                    description = desc_match.group(1)

                fields.append(FieldInfo(
                    name=field_name,
                    type_hint=type_hint,
                    default="None",
                    description=description,
                ))
                i += 1
                continue

            # Match simple field definition: name: Type = value (or just name: Type)
            # This handles: field_name: Optional[StrictStr] = None
            simple_match = re.match(r'^(\w+):\s*(.+?)(?:\s*=\s*(.+))?$', stripped)
            if simple_match and not simple_match.group(1).startswith('_'):
                field_name = simple_match.group(1)
                type_hint = simple_match.group(2).strip()
                default_value = simple_match.group(3).strip() if simple_match.group(3) else None

                # Skip if this looks like a method or class variable annotation
                if type_hint.startswith('ClassVar') or '(' in type_hint:
                    i += 1
                    continue

                # Extract inner type from Optional[]
                optional_match = re.match(r'Optional\[(.+)\]', type_hint)
                if optional_match:
                    type_hint = optional_match.group(1)

                fields.append(FieldInfo(
                    name=field_name,
                    type_hint=type_hint,
                    default=default_value or "None",
                    description="",
                ))

        i += 1

    return fields


def parse_api_methods(api_path: Path) -> list[MethodInfo]:
    """Parse an API file to extract method information."""
    content = api_path.read_text()
    methods = []

    # Find all async def method signatures
    # Pattern: async def method_name(self, params...) -> ReturnType:
    pattern = re.compile(
        r'async def (\w+)\(\s*self,([^)]*)\)\s*->\s*([^:]+):.*?"""(.+?)"""',
        re.DOTALL
    )

    for match in pattern.finditer(content):
        method_name = match.group(1)

        # Skip internal methods and variants
        if (method_name.startswith('_') or
            '_with_http_info' in method_name or
            '_without_preload_content' in method_name):
            continue

        params_str = match.group(2)
        return_type = match.group(3).strip()
        docstring = match.group(4).strip().split('\n')[0]  # First line only

        # Parse parameters
        path_params = []
        request_type = None

        # Split params by comma, but handle nested brackets
        param_parts = []
        current = ""
        bracket_depth = 0
        for char in params_str:
            if char in '([{':
                bracket_depth += 1
            elif char in ')]}':
                bracket_depth -= 1
            elif char == ',' and bracket_depth == 0:
                param_parts.append(current.strip())
                current = ""
                continue
            current += char
        if current.strip():
            param_parts.append(current.strip())

        for param in param_parts:
            param = param.strip()
            if not param or param.startswith('_'):
                continue

            # Parse: name: Type or name: Type = default
            param_match = re.match(r'(\w+):\s*(.+?)(?:\s*=\s*(.+))?$', param)
            if param_match:
                p_name = param_match.group(1)
                p_type = param_match.group(2).strip()
                p_default = param_match.group(3).strip() if param_match.group(3) else None

                # Check if this is a request body parameter with generic 'object' type
                # This happens when OpenAPI spec has requestBody with just 'type: object'
                # Mark it as a special request type that needs no fields
                if p_name == 'request' and p_type == 'object':
                    request_type = '__empty_object__'
                    continue

                # Check if it's a model parameter (request/query type)
                # Look for capitalized type names that aren't basic types
                type_match = re.search(r'(?:Optional\[)?([A-Z]\w+?)(?:\])?$', p_type)
                if type_match:
                    type_name = type_match.group(1)
                    # Skip basic types
                    if type_name not in ('Dict', 'List', 'Optional', 'Union', 'Any', 'Tuple'):
                        if 'Request' in type_name or 'Query' in type_name:
                            request_type = type_name
                        else:
                            # It's a path/query param with a model type
                            path_params.append((p_name, p_type, p_default))
                    else:
                        path_params.append((p_name, p_type, p_default))
                else:
                    path_params.append((p_name, p_type, p_default))

        methods.append(MethodInfo(
            name=method_name,
            path_params=path_params,
            request_type=request_type,
            return_type=return_type,
            docstring=docstring,
        ))

    return methods


def simplify_type(type_str: str) -> str:
    """Convert Pydantic strict types to simple Python types."""
    type_str = type_str.strip()
    replacements = {
        'StrictStr': 'str',
        'StrictInt': 'int',
        'StrictFloat': 'float',
        'StrictBool': 'bool',
        'StrictBytes': 'bytes',
    }
    for old, new in replacements.items():
        type_str = type_str.replace(old, new)
    return type_str


def discover_apis(sdk_dir: Path) -> list[dict]:
    """Discover all API classes and their methods."""
    api_dir = sdk_dir / "api"
    apis = []

    for api_file in sorted(api_dir.glob("*_api.py")):
        module_name = api_file.stem
        content = api_file.read_text()

        class_match = re.search(r'class (\w+Api)\s*:', content)
        if class_match:
            class_name = class_match.group(1)
            property_name = module_name.replace('_api', '')
            methods = parse_api_methods(api_file)

            apis.append({
                'module': module_name,
                'class_name': class_name,
                'property_name': property_name,
                'methods': methods,
            })

    return apis


def discover_models(sdk_dir: Path) -> dict[str, dict]:
    """Discover all models (BaseModel classes and Enums) and their fields."""
    models_dir = sdk_dir / "models"
    models = {}

    for model_file in models_dir.glob("*.py"):
        content = model_file.read_text()

        # Find any class that extends BaseModel
        class_match = re.search(r'class (\w+)\(BaseModel\):', content)
        if class_match:
            class_name = class_match.group(1)
            fields = parse_model_fields(model_file) if 'Request' in class_name else []
            models[class_name] = {
                'fields': fields,
                'module': model_file.stem,
            }
            continue

        # Also find Enum classes (e.g., class SomeStatus(str, Enum):)
        enum_match = re.search(r'class (\w+)\([^)]*Enum[^)]*\):', content)
        if enum_match:
            class_name = enum_match.group(1)
            models[class_name] = {
                'fields': [],
                'module': model_file.stem,
            }

    return models


def generate_wrapper_method(method: MethodInfo, models: dict) -> str:
    """Generate a wrapper method with flattened parameters."""
    lines = []

    # Get request fields if applicable
    request_fields = []
    if method.request_type and method.request_type in models:
        request_fields = models[method.request_type]['fields']

    # Build parameter list
    all_params = []

    # Path params first (required)
    for p_name, p_type, p_default in method.path_params:
        if p_default:
            all_params.append(f"{p_name}: {p_type} = {p_default}")
        else:
            all_params.append(f"{p_name}: {p_type}")

    # Then request fields (all optional)
    for field in request_fields:
        field_type = simplify_type(field.type_hint)
        all_params.append(f"{field.name}: Optional[{field_type}] = None")

    # Method signature
    if all_params:
        params_str = ",\n        ".join(all_params)
        lines.append(f"    async def {method.name}(")
        lines.append(f"        self,")
        lines.append(f"        {params_str},")
        lines.append(f"    ) -> {method.return_type}:")
    else:
        lines.append(f"    async def {method.name}(self) -> {method.return_type}:")

    # Docstring
    lines.append(f'        """{method.docstring}')
    if all_params:
        lines.append("")
        lines.append("        Args:")
        for p_name, p_type, _ in method.path_params:
            lines.append(f"            {p_name}: {p_type}")
        for field in request_fields:
            desc = field.description or field.name
            lines.append(f"            {field.name}: {desc}")
    lines.append('        """')

    # Method body - determine if we need to pass a request object
    call_args = [f"{p[0]}={p[0]}" for p in method.path_params]

    if method.request_type:
        # We have a request type - need to build and pass it
        if method.request_type == '__empty_object__':
            # Special case: request body is just 'object' type with no schema
            # Pass an empty dict as the request
            call_args.append("request={}")
        elif request_fields:
            # Build request object with fields
            lines.append(f"        request = {method.request_type}(")
            for field in request_fields:
                lines.append(f"            {field.name}={field.name},")
            lines.append("        )")
            call_args.append("request=request")
        else:
            # Request type exists but no fields found - create empty request
            lines.append(f"        request = {method.request_type}()")
            call_args.append("request=request")

        lines.append(f"        return await self._api.{method.name}({', '.join(call_args)})")
    else:
        # No request object needed
        if call_args:
            lines.append(f"        return await self._api.{method.name}({', '.join(call_args)})")
        else:
            lines.append(f"        return await self._api.{method.name}()")

    lines.append("")
    return "\n".join(lines)


def generate_unified_client(sdk_dir: Path, package_name: str = "virsh_sandbox"):
    """Generate the unified VirshSandbox client wrapper with flattened parameters."""

    apis = discover_apis(sdk_dir)
    models = discover_models(sdk_dir)

    # Which APIs use tmux host
    tmux_api_properties = {"command", "file", "tmux", "audit", "health", "human", "plan"}

    # Collect imports
    api_imports = []
    model_imports = set()

    for api in apis:
        api_imports.append(f"from {package_name}.api.{api['module']} import {api['class_name']}")
        for method in api['methods']:
            # Import request/query types
            if method.request_type and method.request_type in models:
                model_info = models[method.request_type]
                model_imports.add(
                    f"from {package_name}.models.{model_info['module']} import {method.request_type}"
                )

                # Also import types used in the request model fields (e.g., enum types)
                for field in model_info.get('fields', []):
                    field_type = field.type_hint
                    # Extract type names from the field type (handles List[], Optional[], etc.)
                    field_type_names = re.findall(r'([A-Z][a-zA-Z0-9_]+)', field_type)
                    for field_type_name in field_type_names:
                        if field_type_name in models and field_type_name not in ('Dict', 'List', 'Optional', 'Union', 'Any', 'Tuple', 'StrictStr', 'StrictInt', 'StrictFloat', 'StrictBool', 'StrictBytes'):
                            field_model_info = models[field_type_name]
                            model_imports.add(
                                f"from {package_name}.models.{field_model_info['module']} import {field_type_name}"
                            )

            # Import return types (if they're in our models dict)
            # Handle wrapped types like List[SomeType], Optional[SomeType], etc.
            return_type = method.return_type.strip()
            # Extract type names from the return type (handles List[], Optional[], etc.)
            type_names = re.findall(r'([A-Z][a-zA-Z0-9_]+)', return_type)
            for type_name in type_names:
                if type_name in models and type_name not in ('Dict', 'List', 'Optional', 'Union', 'Any', 'Tuple'):
                    model_info = models[type_name]
                    model_imports.add(
                        f"from {package_name}.models.{model_info['module']} import {type_name}"
                    )

            # Import types from path params that are model types
            for p_name, p_type, _ in method.path_params:
                type_match = re.search(r'(?:Optional\[)?([A-Z]\w+?)(?:\])?$', p_type)
                if type_match:
                    type_name = type_match.group(1)
                    if type_name in models:
                        model_info = models[type_name]
                        model_imports.add(
                            f"from {package_name}.models.{model_info['module']} import {type_name}"
                        )

    # Generate wrapper classes
    wrapper_classes = []
    for api in apis:
        wrapper_name = api['class_name'].replace('Api', 'Operations')
        lines = []
        lines.append(f"class {wrapper_name}:")
        lines.append(f'    """Wrapper for {api["class_name"]} with simplified method signatures."""')
        lines.append("")
        lines.append(f"    def __init__(self, api: {api['class_name']}):")
        lines.append("        self._api = api")
        lines.append("")

        for method in api['methods']:
            method_code = generate_wrapper_method(method, models)
            lines.append(method_code)

        wrapper_classes.append("\n".join(lines))

    # Build the complete file
    output_lines = []

    output_lines.append('# coding: utf-8')
    output_lines.append('')
    output_lines.append('"""')
    output_lines.append('Unified VirshSandbox Client')
    output_lines.append('')
    output_lines.append('This module provides a unified client wrapper for the virsh-sandbox SDK,')
    output_lines.append('offering a cleaner interface with flattened parameters instead of request objects.')
    output_lines.append('')
    output_lines.append('Example:')
    output_lines.append(f'    from {package_name} import VirshSandbox')
    output_lines.append('')
    output_lines.append('    async with VirshSandbox(host="http://localhost:8080") as client:')
    output_lines.append('        # Create a sandbox with simple parameters')
    output_lines.append('        await client.sandbox.create_sandbox(source_vm_name="ubuntu-base")')
    output_lines.append('        # Run a command')
    output_lines.append('        await client.command.run_command(command="ls", args=["-la"])')
    output_lines.append('"""')
    output_lines.append('')
    output_lines.append('from typing import Any, Dict, List, Optional, Tuple, Union')
    output_lines.append('')
    output_lines.append(f'from {package_name}.api_client import ApiClient')
    output_lines.append(f'from {package_name}.configuration import Configuration')

    for imp in sorted(api_imports):
        output_lines.append(imp)

    for imp in sorted(model_imports):
        output_lines.append(imp)

    output_lines.append('')
    output_lines.append('')

    # Add wrapper classes
    for wrapper in wrapper_classes:
        output_lines.append(wrapper)
        output_lines.append('')

    # Main client class
    output_lines.append('')
    output_lines.append('class VirshSandbox:')
    output_lines.append('    """Unified client for the virsh-sandbox API.')
    output_lines.append('')
    output_lines.append('    This class provides a single entry point for all virsh-sandbox API operations,')
    output_lines.append('    with support for separate hosts for the main API and tmux API.')
    output_lines.append('    All methods use flattened parameters instead of request objects.')
    output_lines.append('')
    output_lines.append('    Args:')
    output_lines.append('        host: Base URL for the main virsh-sandbox API')
    output_lines.append('        tmux_host: Base URL for the tmux API (defaults to host)')
    output_lines.append('        api_key: Optional API key for authentication')
    output_lines.append('        verify_ssl: Whether to verify SSL certificates')
    output_lines.append('')
    output_lines.append('    Example:')
    output_lines.append(f'        >>> from {package_name} import VirshSandbox')
    output_lines.append('        >>> async with VirshSandbox() as client:')
    output_lines.append('        ...     await client.sandbox.create_sandbox(source_vm_name="base-vm")')
    output_lines.append('    """')
    output_lines.append('')
    output_lines.append('    def __init__(')
    output_lines.append('        self,')
    output_lines.append('        host: str = "http://localhost:8080",')
    output_lines.append('        tmux_host: Optional[str] = None,')
    output_lines.append('        api_key: Optional[str] = None,')
    output_lines.append('        access_token: Optional[str] = None,')
    output_lines.append('        username: Optional[str] = None,')
    output_lines.append('        password: Optional[str] = None,')
    output_lines.append('        verify_ssl: bool = True,')
    output_lines.append('        ssl_ca_cert: Optional[str] = None,')
    output_lines.append('        retries: Optional[int] = None,')
    output_lines.append('    ) -> None:')
    output_lines.append('        """Initialize the VirshSandbox client."""')
    output_lines.append('        self._main_config = Configuration(')
    output_lines.append('            host=host,')
    output_lines.append('            api_key={"Authorization": api_key} if api_key else None,')
    output_lines.append('            access_token=access_token,')
    output_lines.append('            username=username,')
    output_lines.append('            password=password,')
    output_lines.append('            ssl_ca_cert=ssl_ca_cert,')
    output_lines.append('            retries=retries,')
    output_lines.append('        )')
    output_lines.append('        self._main_config.verify_ssl = verify_ssl')
    output_lines.append('        self._main_api_client = ApiClient(configuration=self._main_config)')
    output_lines.append('')
    output_lines.append('        tmux_host = tmux_host or host')
    output_lines.append('        if tmux_host != host:')
    output_lines.append('            self._tmux_config = Configuration(')
    output_lines.append('                host=tmux_host,')
    output_lines.append('                api_key={"Authorization": api_key} if api_key else None,')
    output_lines.append('                access_token=access_token,')
    output_lines.append('                username=username,')
    output_lines.append('                password=password,')
    output_lines.append('                ssl_ca_cert=ssl_ca_cert,')
    output_lines.append('                retries=retries,')
    output_lines.append('            )')
    output_lines.append('            self._tmux_config.verify_ssl = verify_ssl')
    output_lines.append('            self._tmux_api_client = ApiClient(configuration=self._tmux_config)')
    output_lines.append('        else:')
    output_lines.append('            self._tmux_config = self._main_config')
    output_lines.append('            self._tmux_api_client = self._main_api_client')
    output_lines.append('')

    # Lazy init fields
    for api in apis:
        wrapper_name = api['class_name'].replace('Api', 'Operations')
        output_lines.append(f"        self._{api['property_name']}: Optional[{wrapper_name}] = None")
    output_lines.append('')

    # Properties
    for api in apis:
        prop = api['property_name']
        wrapper_name = api['class_name'].replace('Api', 'Operations')
        api_class = api['class_name']
        client_var = "self._tmux_api_client" if prop in tmux_api_properties else "self._main_api_client"

        output_lines.append('    @property')
        output_lines.append(f'    def {prop}(self) -> {wrapper_name}:')
        output_lines.append(f'        """Access {api_class} operations."""')
        output_lines.append(f'        if self._{prop} is None:')
        output_lines.append(f'            api = {api_class}(api_client={client_var})')
        output_lines.append(f'            self._{prop} = {wrapper_name}(api)')
        output_lines.append(f'        return self._{prop}')
        output_lines.append('')

    # Utility methods
    output_lines.append('    @property')
    output_lines.append('    def configuration(self) -> Configuration:')
    output_lines.append('        """Get the main API configuration."""')
    output_lines.append('        return self._main_config')
    output_lines.append('')
    output_lines.append('    @property')
    output_lines.append('    def tmux_configuration(self) -> Configuration:')
    output_lines.append('        """Get the tmux API configuration."""')
    output_lines.append('        return self._tmux_config')
    output_lines.append('')
    output_lines.append('    def set_debug(self, debug: bool) -> None:')
    output_lines.append('        """Enable or disable debug mode."""')
    output_lines.append('        self._main_config.debug = debug')
    output_lines.append('        if self._tmux_config is not self._main_config:')
    output_lines.append('            self._tmux_config.debug = debug')
    output_lines.append('')
    output_lines.append('    async def close(self) -> None:')
    output_lines.append('        """Close the API client connections."""')
    output_lines.append("        if hasattr(self._main_api_client.rest_client, 'close'):")
    output_lines.append('            await self._main_api_client.rest_client.close()')
    output_lines.append('        if self._tmux_api_client is not self._main_api_client:')
    output_lines.append("            if hasattr(self._tmux_api_client.rest_client, 'close'):")
    output_lines.append('                await self._tmux_api_client.rest_client.close()')
    output_lines.append('')
    output_lines.append('    async def __aenter__(self) -> "VirshSandbox":')
    output_lines.append('        """Async context manager entry."""')
    output_lines.append('        return self')
    output_lines.append('')
    output_lines.append('    async def __aexit__(self, exc_type, exc_val, exc_tb) -> None:')
    output_lines.append('        """Async context manager exit."""')
    output_lines.append('        await self.close()')

    # Write the file
    client_path = sdk_dir / "client.py"
    client_path.write_text("\n".join(output_lines))
    print(f"Generated unified client: {client_path}")


def update_init_file(sdk_dir: Path, package_name: str = "virsh_sandbox"):
    """Update __init__.py to export VirshSandbox."""
    init_path = sdk_dir / "__init__.py"
    content = init_path.read_text()

    if f"from {package_name}.client import VirshSandbox" in content:
        print("VirshSandbox already exported in __init__.py")
        return

    content = content.replace(
        '__all__ = [',
        '__all__ = [\n    "VirshSandbox",'
    )

    if "# import apis into sdk package" in content:
        content = content.replace(
            "# import apis into sdk package",
            f"# import unified client\nfrom {package_name}.client import VirshSandbox as VirshSandbox\n\n# import apis into sdk package"
        )
    else:
        content += f"\n# import unified client\nfrom {package_name}.client import VirshSandbox as VirshSandbox\n"

    init_path.write_text(content)
    print("Updated __init__.py to export VirshSandbox")


def main():
    sdk_dir = Path("virsh-sandbox-py/virsh_sandbox")
    package_name = "virsh_sandbox"

    if not sdk_dir.exists():
        print(f"SDK directory not found: {sdk_dir}")
        print("Make sure to run this script from the sdk/ directory")
        return

    print("Generating unified client with flattened parameters...")
    generate_unified_client(sdk_dir, package_name)

    print("Updating __init__.py...")
    update_init_file(sdk_dir, package_name)

    print("SDK polished!")


if __name__ == "__main__":
    main()
