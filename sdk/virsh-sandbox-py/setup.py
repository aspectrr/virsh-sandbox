from setuptools import find_packages, setup

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="virsh_sandbox",
    version="0.0.5-beta",
    author="Collin Pfeifer",
    author_email="cpfeifer@madcactus.org",
    description="API for managing virtual machine sandboxes using libvirt",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/aspectrr/virsh-sandbox",
    packages=find_packages(exclude=["tests", "tests.*"]),
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
    ],
    python_requires=">=3.8",
    install_requires=[
        "urllib3>=1.25.3",
        "python-dateutil",
        "pydantic>=2.0.0",
        "typing-extensions>=4.0.0",
    ],
    extras_require={
        "dev": [
            "pytest>=7.0.0",
            "pytest-cov",
            "black",
            "isort",
            "mypy",
            "flake8",
        ],
    },
    package_data={
        "virsh_sandbox": ["py.typed"],
    },
)
