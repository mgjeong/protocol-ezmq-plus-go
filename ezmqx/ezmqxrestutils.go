/*******************************************************************************
 * Copyright 2018 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/

package ezmqx

// URLs
const NODE = "http://pharos-node:48098"
const PREFIX = "/api/v1"
const API_CONFIG = "/management/device/configuration"
const API_APPS = "/management/apps"
const API_DETAIL = "/management/detail"
const TNS_KNOWN_PORT = "48323"
const TOPIC = "/tns/topic"
const TNS_KEEP_ALIVE = "/tns/keepalive"
const HTTP_PREFIX = "http://"
const QUERY_NAME = "name="
const QUERY_HIERARCHICAL = "&hierarchical="
const QUERY_TRUE = "yes"
const QUERY_FALSE = "no"
const REVERSE_PROXY_KNOWN_PORT = "80"
const API_SEARCH_NODE = "/search/nodes"
const REVERSE_PROXY_PREFIX = "/tns-server"
const ANCHOR_IMAGE_NAME = "imageName="

// JSON Keys
const CONF_PROPS = "properties"
const CONF_ANCHOR_ADDR = "anchorendpoint"
const CONF_NODE_ADDR = "nodeaddress"
const APPS_PROPS = "apps"
const APPS_ID = "id"
const APPS_STATE = "state"
const APPS_STATE_RUNNING = "running"
const APPS_EXIT_CODE = "exitcode"
const APPS_STATUS = "status"
const SERVICES_PROPS = "services"
const SERVICES_CON_NAME = "name"
const SERVICES_CON_ID = "cid"
const SERVICES_CON_PORTS = "ports"
const PORTS_PRIVATE = "PrivatePort"
const PORTS_PUBLIC = "PublicPort"
const PAYLOAD_OPTION = "indentation"
const PAYLOAD_TOPIC = "topic"
const PAYLOAD_TOPICS = "topics"
const PAYLOAD_NAME = "name"
const PAYLOAD_ENDPOINT = "endpoint"
const PAYLOAD_DATAMODEL = "datamodel"
const PAYLOAD_KEEPALIVE_INTERVAL = "ka_interval"
const PAYLOAD_TOPIC_KA = "topic_names"
const CONF_REVERSE_PROXY = "reverseproxy"
const CONF_REVERSE_PROXY_ENABLED = "enabled"
const NODES = "nodes"
const NODES_STATUS = "status"
const NODES_CONNECTED = "connected"
const NODES_IP = "ip"
const NODES_CONF = "config"
const NODES_PROPS = "properties"
const NODES_REVERSE_PROXY = "reverseproxy"
const NODES_REVERSE_PROXY_ENABLED = "enabled"
const CONFIG_ANCHOR_IMAGE_NAME = "imageName"

// HostName file path
const HOST_NAME_FILE_PATH = "/etc/hostname"

// HTTP status codes
const HTTP_OK = 200
const HTTP_CREATED = 201
const CONNECTION_TIMEOUT = 5

// Strings
const SLASH = "/"
const DOUBLE_SLASH = "//"
const COLON = ":"
const QUESTION_MARK = "?"
const REGISTER = "register"
const UNREGISTER = "unregister"
const KEEPALIVE = "keepalive"
const SHUTDOWN = "shutdown"
const APPLICATION_JSON = "application/json"
