# coding: utf-8

"""
    virsh-sandbox API
    API for managing virtual machine sandboxes using libvirt
"""

from typing import Any, Dict, List, Optional, Tuple, Union

from pydantic import Field, StrictStr
from typing_extensions import Annotated

from virsh_sandbox.api_client import ApiClient, RequestSerialized
from virsh_sandbox.api_response import ApiResponse
from virsh_sandbox.exceptions import ApiException
from virsh_sandbox.models.tmux_client_internal_types_create_pane_request import \
    TmuxClientInternalTypesCreatePaneRequest
from virsh_sandbox.models.tmux_client_internal_types_create_pane_response import \
    TmuxClientInternalTypesCreatePaneResponse
from virsh_sandbox.models.tmux_client_internal_types_kill_session_response import \
    TmuxClientInternalTypesKillSessionResponse
from virsh_sandbox.models.tmux_client_internal_types_list_panes_response import \
    TmuxClientInternalTypesListPanesResponse
from virsh_sandbox.models.tmux_client_internal_types_read_pane_request import \
    TmuxClientInternalTypesReadPaneRequest
from virsh_sandbox.models.tmux_client_internal_types_read_pane_response import \
    TmuxClientInternalTypesReadPaneResponse
from virsh_sandbox.models.tmux_client_internal_types_send_keys_request import \
    TmuxClientInternalTypesSendKeysRequest
from virsh_sandbox.models.tmux_client_internal_types_send_keys_response import \
    TmuxClientInternalTypesSendKeysResponse
from virsh_sandbox.models.tmux_client_internal_types_session_info import \
    TmuxClientInternalTypesSessionInfo
from virsh_sandbox.models.tmux_client_internal_types_switch_pane_request import \
    TmuxClientInternalTypesSwitchPaneRequest
from virsh_sandbox.models.tmux_client_internal_types_switch_pane_response import \
    TmuxClientInternalTypesSwitchPaneResponse
from virsh_sandbox.models.tmux_client_internal_types_window_info import \
    TmuxClientInternalTypesWindowInfo


class TmuxApi:
    """TmuxApi service"""

    def __init__(self, api_client: Optional[ApiClient] = None) -> None:
        if api_client is None:
            api_client = ApiClient.get_default()
        self.api_client = api_client

    def create_tmux_pane(
        self,
        request: TmuxClientInternalTypesCreatePaneRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> TmuxClientInternalTypesCreatePaneResponse:
        """Create tmux pane

        Creates a new tmux pane

        :param request: Create pane request (required)
        :type request: TmuxClientInternalTypesCreatePaneRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._create_tmux_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesCreatePaneResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def create_tmux_pane_with_http_info(
        self,
        request: TmuxClientInternalTypesCreatePaneRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[TmuxClientInternalTypesCreatePaneResponse]:
        """Create tmux pane

        Creates a new tmux pane

        :param request: Create pane request (required)
        :type request: TmuxClientInternalTypesCreatePaneRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._create_tmux_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesCreatePaneResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def create_tmux_pane_without_preload_content(
        self,
        request: TmuxClientInternalTypesCreatePaneRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """Create tmux pane

        Creates a new tmux pane

        :param request: Create pane request (required)
        :type request: TmuxClientInternalTypesCreatePaneRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._create_tmux_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesCreatePaneResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _create_tmux_pane_serialize(
        self,
        request: TmuxClientInternalTypesCreatePaneRequest,
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        # process the query parameters
        # process the header parameters
        # process the form parameters
        # process the body parameter
        if request is not None:
            _body_params = request

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`
        if _content_type:
            _header_params["Content-Type"] = _content_type
        else:
            _default_content_type = self.api_client.select_header_content_type(
                ["application/json"]
            )
            if _default_content_type is not None:
                _header_params["Content-Type"] = _default_content_type

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="POST",
            resource_path="/v1/tmux/panes/create",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def create_tmux_session(
        self,
        request: object,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Dict[str, str]:
        """Create tmux session

        Creates a new tmux session

        :param request: Create session request (required)
        :type request: object
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._create_tmux_session_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "Dict[str, str]",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def create_tmux_session_with_http_info(
        self,
        request: object,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[Dict[str, str]]:
        """Create tmux session

        Creates a new tmux session

        :param request: Create session request (required)
        :type request: object
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._create_tmux_session_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "Dict[str, str]",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def create_tmux_session_without_preload_content(
        self,
        request: object,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """Create tmux session

        Creates a new tmux session

        :param request: Create session request (required)
        :type request: object
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._create_tmux_session_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "Dict[str, str]",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _create_tmux_session_serialize(
        self,
        request: object,
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        # process the query parameters
        # process the header parameters
        # process the form parameters
        # process the body parameter
        if request is not None:
            _body_params = request

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`
        if _content_type:
            _header_params["Content-Type"] = _content_type
        else:
            _default_content_type = self.api_client.select_header_content_type(
                ["application/json"]
            )
            if _default_content_type is not None:
                _header_params["Content-Type"] = _default_content_type

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="POST",
            resource_path="/v1/tmux/sessions/create",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def kill_tmux_pane(
        self,
        pane_id: str,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Dict[str, object]:
        """Kill tmux pane

        Kills a tmux pane

        :param pane_id: Pane ID (required)
        :type pane_id: str
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._kill_tmux_pane_serialize(
            pane_id=pane_id,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "Dict[str, object]",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def kill_tmux_pane_with_http_info(
        self,
        pane_id: str,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[Dict[str, object]]:
        """Kill tmux pane

        Kills a tmux pane

        :param pane_id: Pane ID (required)
        :type pane_id: str
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._kill_tmux_pane_serialize(
            pane_id=pane_id,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "Dict[str, object]",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def kill_tmux_pane_without_preload_content(
        self,
        pane_id: str,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """Kill tmux pane

        Kills a tmux pane

        :param pane_id: Pane ID (required)
        :type pane_id: str
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._kill_tmux_pane_serialize(
            pane_id=pane_id,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "Dict[str, object]",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _kill_tmux_pane_serialize(
        self,
        pane_id: str,
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        if pane_id is not None:
            _path_params["paneID"] = pane_id
        # process the query parameters
        # process the header parameters
        # process the form parameters
        # process the body parameter

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="DELETE",
            resource_path="/v1/tmux/panes/{paneID}",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def kill_tmux_session(
        self,
        session_name: str,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Dict[str, object]:
        """Kill tmux session

        Kills a tmux session

        :param session_name: Session name (required)
        :type session_name: str
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._kill_tmux_session_serialize(
            session_name=session_name,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "Dict[str, object]",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def kill_tmux_session_with_http_info(
        self,
        session_name: str,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[Dict[str, object]]:
        """Kill tmux session

        Kills a tmux session

        :param session_name: Session name (required)
        :type session_name: str
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._kill_tmux_session_serialize(
            session_name=session_name,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "Dict[str, object]",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def kill_tmux_session_without_preload_content(
        self,
        session_name: str,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """Kill tmux session

        Kills a tmux session

        :param session_name: Session name (required)
        :type session_name: str
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._kill_tmux_session_serialize(
            session_name=session_name,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "Dict[str, object]",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _kill_tmux_session_serialize(
        self,
        session_name: str,
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        if session_name is not None:
            _path_params["sessionName"] = session_name
        # process the query parameters
        # process the header parameters
        # process the form parameters
        # process the body parameter

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="DELETE",
            resource_path="/v1/tmux/sessions/{sessionName}",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def list_tmux_panes(
        self,
        session: Optional[str] = None,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> TmuxClientInternalTypesListPanesResponse:
        """List tmux panes

        Get a list of panes in a tmux session

        :param session: Session name (optional)
        :type session: str, optional
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._list_tmux_panes_serialize(
            session=session,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesListPanesResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def list_tmux_panes_with_http_info(
        self,
        session: Optional[str] = None,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[TmuxClientInternalTypesListPanesResponse]:
        """List tmux panes

        Get a list of panes in a tmux session

        :param session: Session name (optional)
        :type session: str, optional
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._list_tmux_panes_serialize(
            session=session,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesListPanesResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def list_tmux_panes_without_preload_content(
        self,
        session: Optional[str] = None,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """List tmux panes

        Get a list of panes in a tmux session

        :param session: Session name (optional)
        :type session: str, optional
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._list_tmux_panes_serialize(
            session=session,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesListPanesResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _list_tmux_panes_serialize(
        self,
        session: Optional[str],
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        # process the query parameters
        if session is not None:
            _query_params.append(("session", session))

        # process the header parameters
        # process the form parameters
        # process the body parameter

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="GET",
            resource_path="/v1/tmux/panes",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def list_tmux_sessions(
        self,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> List[TmuxClientInternalTypesSessionInfo]:
        """List tmux sessions

        Get a list of all active tmux sessions

        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._list_tmux_sessions_serialize(
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "List[TmuxClientInternalTypesSessionInfo]",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def list_tmux_sessions_with_http_info(
        self,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[List[TmuxClientInternalTypesSessionInfo]]:
        """List tmux sessions

        Get a list of all active tmux sessions

        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._list_tmux_sessions_serialize(
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "List[TmuxClientInternalTypesSessionInfo]",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def list_tmux_sessions_without_preload_content(
        self,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """List tmux sessions

        Get a list of all active tmux sessions

        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._list_tmux_sessions_serialize(
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "List[TmuxClientInternalTypesSessionInfo]",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _list_tmux_sessions_serialize(
        self,
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        # process the query parameters
        # process the header parameters
        # process the form parameters
        # process the body parameter

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="GET",
            resource_path="/v1/tmux/sessions",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def list_tmux_windows(
        self,
        session: Optional[str] = None,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> List[TmuxClientInternalTypesWindowInfo]:
        """List tmux windows

        Get a list of windows in a tmux session

        :param session: Session name (optional)
        :type session: str, optional
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._list_tmux_windows_serialize(
            session=session,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "List[TmuxClientInternalTypesWindowInfo]",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def list_tmux_windows_with_http_info(
        self,
        session: Optional[str] = None,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[List[TmuxClientInternalTypesWindowInfo]]:
        """List tmux windows

        Get a list of windows in a tmux session

        :param session: Session name (optional)
        :type session: str, optional
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._list_tmux_windows_serialize(
            session=session,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "List[TmuxClientInternalTypesWindowInfo]",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def list_tmux_windows_without_preload_content(
        self,
        session: Optional[str] = None,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """List tmux windows

        Get a list of windows in a tmux session

        :param session: Session name (optional)
        :type session: str, optional
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._list_tmux_windows_serialize(
            session=session,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "List[TmuxClientInternalTypesWindowInfo]",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _list_tmux_windows_serialize(
        self,
        session: Optional[str],
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        # process the query parameters
        if session is not None:
            _query_params.append(("session", session))

        # process the header parameters
        # process the form parameters
        # process the body parameter

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="GET",
            resource_path="/v1/tmux/windows",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def read_tmux_pane(
        self,
        request: TmuxClientInternalTypesReadPaneRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> TmuxClientInternalTypesReadPaneResponse:
        """Read tmux pane

        Reads the content of a tmux pane

        :param request: Read pane request (required)
        :type request: TmuxClientInternalTypesReadPaneRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._read_tmux_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesReadPaneResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def read_tmux_pane_with_http_info(
        self,
        request: TmuxClientInternalTypesReadPaneRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[TmuxClientInternalTypesReadPaneResponse]:
        """Read tmux pane

        Reads the content of a tmux pane

        :param request: Read pane request (required)
        :type request: TmuxClientInternalTypesReadPaneRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._read_tmux_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesReadPaneResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def read_tmux_pane_without_preload_content(
        self,
        request: TmuxClientInternalTypesReadPaneRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """Read tmux pane

        Reads the content of a tmux pane

        :param request: Read pane request (required)
        :type request: TmuxClientInternalTypesReadPaneRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._read_tmux_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesReadPaneResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _read_tmux_pane_serialize(
        self,
        request: TmuxClientInternalTypesReadPaneRequest,
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        # process the query parameters
        # process the header parameters
        # process the form parameters
        # process the body parameter
        if request is not None:
            _body_params = request

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`
        if _content_type:
            _header_params["Content-Type"] = _content_type
        else:
            _default_content_type = self.api_client.select_header_content_type(
                ["application/json"]
            )
            if _default_content_type is not None:
                _header_params["Content-Type"] = _default_content_type

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="POST",
            resource_path="/v1/tmux/panes/read",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def release_tmux_session(
        self,
        session_id: str,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> TmuxClientInternalTypesKillSessionResponse:
        """Release tmux session

        Releases (kills) a tmux session by ID

        :param session_id: Session ID (required)
        :type session_id: str
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._release_tmux_session_serialize(
            session_id=session_id,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesKillSessionResponse",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def release_tmux_session_with_http_info(
        self,
        session_id: str,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[TmuxClientInternalTypesKillSessionResponse]:
        """Release tmux session

        Releases (kills) a tmux session by ID

        :param session_id: Session ID (required)
        :type session_id: str
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._release_tmux_session_serialize(
            session_id=session_id,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesKillSessionResponse",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def release_tmux_session_without_preload_content(
        self,
        session_id: str,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """Release tmux session

        Releases (kills) a tmux session by ID

        :param session_id: Session ID (required)
        :type session_id: str
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._release_tmux_session_serialize(
            session_id=session_id,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesKillSessionResponse",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _release_tmux_session_serialize(
        self,
        session_id: str,
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        if session_id is not None:
            _path_params["sessionId"] = session_id
        # process the query parameters
        # process the header parameters
        # process the form parameters
        # process the body parameter

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="POST",
            resource_path="/v1/tmux/sessions/{sessionId}/release",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def send_keys_to_pane(
        self,
        request: TmuxClientInternalTypesSendKeysRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> TmuxClientInternalTypesSendKeysResponse:
        """Send keys to tmux pane

        Sends keystrokes to a tmux pane

        :param request: Send keys request (required)
        :type request: TmuxClientInternalTypesSendKeysRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._send_keys_to_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesSendKeysResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "403": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def send_keys_to_pane_with_http_info(
        self,
        request: TmuxClientInternalTypesSendKeysRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[TmuxClientInternalTypesSendKeysResponse]:
        """Send keys to tmux pane

        Sends keystrokes to a tmux pane

        :param request: Send keys request (required)
        :type request: TmuxClientInternalTypesSendKeysRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._send_keys_to_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesSendKeysResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "403": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def send_keys_to_pane_without_preload_content(
        self,
        request: TmuxClientInternalTypesSendKeysRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """Send keys to tmux pane

        Sends keystrokes to a tmux pane

        :param request: Send keys request (required)
        :type request: TmuxClientInternalTypesSendKeysRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._send_keys_to_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesSendKeysResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "403": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _send_keys_to_pane_serialize(
        self,
        request: TmuxClientInternalTypesSendKeysRequest,
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        # process the query parameters
        # process the header parameters
        # process the form parameters
        # process the body parameter
        if request is not None:
            _body_params = request

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`
        if _content_type:
            _header_params["Content-Type"] = _content_type
        else:
            _default_content_type = self.api_client.select_header_content_type(
                ["application/json"]
            )
            if _default_content_type is not None:
                _header_params["Content-Type"] = _default_content_type

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="POST",
            resource_path="/v1/tmux/panes/send-keys",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )

    def switch_tmux_pane(
        self,
        request: TmuxClientInternalTypesSwitchPaneRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> TmuxClientInternalTypesSwitchPaneResponse:
        """Switch tmux pane

        Switches to a specific tmux pane

        :param request: Switch pane request (required)
        :type request: TmuxClientInternalTypesSwitchPaneRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object.
        """

        _param = self._switch_tmux_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesSwitchPaneResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        ).data

    def switch_tmux_pane_with_http_info(
        self,
        request: TmuxClientInternalTypesSwitchPaneRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> ApiResponse[TmuxClientInternalTypesSwitchPaneResponse]:
        """Switch tmux pane

        Switches to a specific tmux pane

        :param request: Switch pane request (required)
        :type request: TmuxClientInternalTypesSwitchPaneRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object with HTTP info.
        """

        _param = self._switch_tmux_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesSwitchPaneResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        response_data.read()
        return self.api_client.response_deserialize(
            response_data=response_data,
            response_types_map=_response_types_map,
        )

    def switch_tmux_pane_without_preload_content(
        self,
        request: TmuxClientInternalTypesSwitchPaneRequest,
        _request_timeout: Union[None, float, Tuple[float, float]] = None,
        _request_auth: Optional[Dict[str, Any]] = None,
        _content_type: Optional[str] = None,
        _headers: Optional[Dict[str, Any]] = None,
        _host_index: int = 0,
    ) -> Any:
        """Switch tmux pane

        Switches to a specific tmux pane

        :param request: Switch pane request (required)
        :type request: TmuxClientInternalTypesSwitchPaneRequest
        :param _request_timeout: Timeout setting for this request. If one
                                 number is provided, it will be the total request
                                 timeout. It can also be a pair (tuple) of
                                 (connection, read) timeouts.
        :type _request_timeout: int, tuple(int, int), optional
        :param _request_auth: Override the auth_settings for a single request.
        :type _request_auth: dict, optional
        :param _content_type: Force content-type for the request.
        :type _content_type: str, optional
        :param _headers: Override headers for a single request.
        :type _headers: dict, optional
        :param _host_index: Override host index for a single request.
        :type _host_index: int, optional
        :return: Returns the result object without preloading content.
        """

        _param = self._switch_tmux_pane_serialize(
            request=request,
            _request_auth=_request_auth,
            _content_type=_content_type,
            _headers=_headers,
            _host_index=_host_index,
        )

        _response_types_map: Dict[str, Optional[str]] = {
            "200": "TmuxClientInternalTypesSwitchPaneResponse",
            "400": "TmuxClientInternalTypesAPIError",
            "500": "TmuxClientInternalTypesAPIError",
        }
        response_data = self.api_client.call_api(
            *_param, _request_timeout=_request_timeout
        )
        return response_data.response

    def _switch_tmux_pane_serialize(
        self,
        request: TmuxClientInternalTypesSwitchPaneRequest,
        _request_auth: Optional[Dict[str, Any]],
        _content_type: Optional[str],
        _headers: Optional[Dict[str, Any]],
        _host_index: int,
    ) -> RequestSerialized:
        _host = None

        _collection_formats: Dict[str, str] = {}

        _path_params: Dict[str, str] = {}
        _query_params: List[Tuple[str, str]] = []
        _header_params: Dict[str, Optional[str]] = _headers or {}
        _form_params: List[Tuple[str, str]] = []
        _files: Dict[
            str, Union[str, bytes, List[str], List[bytes], List[Tuple[str, bytes]]]
        ] = {}
        _body_params: Any = None

        # process the path parameters
        # process the query parameters
        # process the header parameters
        # process the form parameters
        # process the body parameter
        if request is not None:
            _body_params = request

        # set the HTTP header `Accept`
        if "Tmux" not in _header_params:
            _header_params["Accept"] = self.api_client.select_header_accept(
                ["application/json"]
            )

        # set the HTTP header `Content-Type`
        if _content_type:
            _header_params["Content-Type"] = _content_type
        else:
            _default_content_type = self.api_client.select_header_content_type(
                ["application/json"]
            )
            if _default_content_type is not None:
                _header_params["Content-Type"] = _default_content_type

        # authentication setting
        _auth_settings: List[str] = []

        return self.api_client.param_serialize(
            method="POST",
            resource_path="/v1/tmux/panes/switch",
            path_params=_path_params,
            query_params=_query_params,
            header_params=_header_params,
            body=_body_params,
            post_params=_form_params,
            files=_files,
            auth_settings=_auth_settings,
            collection_formats=_collection_formats,
            _host=_host,
            _request_auth=_request_auth,
        )
