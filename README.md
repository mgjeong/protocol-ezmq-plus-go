# ezmq-plus library (go)

protocol-ezmq-plus-go is a go package which provides a standard messaging interface over various data streaming
and serialization / deserialization middlewares along with some added functionalities.</br>
  - Currently supports streaming using 0mq and serialization / deserialization using protobuf.
  - Publisher -> Multiple Subscribers broadcasting.
  - High speed serialization and deserialization.
  - Topic name discovery [TNS]. 
  - Automation Markup language [AML]

## Prerequisites ##
- Go compiler
  - Version : 1.9
  - [How to install](https://golang.org/doc/install)
- protocol-ezmq-go
  - Since [protocol-ezmq-go](https://github.sec.samsung.net/RS7-EdgeComputing/protocol-ezmq-go) will be downloaded and built when protocol-ezmq-plus-go is built, check the prerequisites of it. It can be installed via build option (See 'How to build')
- datamodel-aml-go
  - Since [datamodel-aml-go](https://github.sec.samsung.net/RS7-EdgeComputing/datamodel-aml-go) will be downloaded and built when protocol-ezmq-plus-gois built, check the prerequisites of it. It can be installed via build option (See 'How to build')

## How to build ##
1. Goto: ~/protocol-ezmq-plus-go/
2. Run the script:

```
   ./build.sh         : Native build for x86_64
   ./build_arm.sh     : Native build for armhf [Raspberry Pi]
```
   
**Note:** </br>
1. For getting help about script option: **$ ./build_common.sh --help** </br>
2. If building for first time, use **--with_dependencies=true** flag.

## How to run ezmq-plus samples [binary/executables] ##
ezmq-plus has publisher, amlsubscriber, xmlsubscriber and topic-discovery sample applications. Build and run using the following guide to experiment different options in sample.

### Prerequisites ###
 **For secured sample** : Built ezmq-plus library with secured flag.</br>
 **For unsecured sample** : Built ezmq-plus library without secured flag.</br>

### Publisher sample [Secured] ###
1. Goto: ~/protocol-ezmq-plus-go/src/go/ezmqx_samples
2. export LD_LIBRARY_PATH=../ezmqx_extlibs/
3. Run the sample:
    ```
    $ ./publisher_secured
    ```
**Note:** It will give list of options for running the sample. 

### AML Subscriber sample [Secured] ###
1. Goto: ~/protocol-ezmq-plus-go/src/go/ezmqx_samples
2. export LD_LIBRARY_PATH=../ezmqx_extlibs/
3. Run the sample:
    ```
    $ ./amlsubscriber_secured
    ```
**Note:** It will give list of options for running the sample.  

### XML Subscriber sample [Secured] ###
1. Goto: ~/protocol-ezmq-plus-go/src/go/ezmqx_samples
2. export LD_LIBRARY_PATH=../ezmqx_extlibs/
3. Run the sample:
    ```
    $ ./xmlsubscriber_secured
    ```
**Note:** It will give list of options for running the sample. 
 
### Topic Discovery sample [Secured/Unsecured] ###
1. Goto: ~/protocol-ezmq-plus-go/src/go/ezmqx_samples
2. export LD_LIBRARY_PATH=../ezmqx_extlibs/
3. Run the sample:
    ```
    $ ./topicdiscovery
    ```
**Note:** It will give list of options for running the sample. 

### Publisher sample [UnSecured]  ###
1. Goto: ~/protocol-ezmq-plus-go/src/go/ezmqx_samples
2. export LD_LIBRARY_PATH=../ezmqx_extlibs/
3. Run the sample:
    ```
    $ ./publisher
    ```
**Note:** It will give list of options for running the sample. 

### AML Subscriber sample [UnSecured]  ###
1. Goto: ~/protocol-ezmq-plus-go/src/go/ezmqx_samples
2. export LD_LIBRARY_PATH=../ezmqx_extlibs/
3. Run the sample:
    ```
    $ ./amlsubscriber
    ```
**Note:** It will give list of options for running the sample.  

### XML Subscriber sample [UnSecured]  ###
1. Goto: ~/protocol-ezmq-plus-go/src/go/ezmqx_samples
2. export LD_LIBRARY_PATH=../ezmqx_extlibs/
3. Run the sample:
    ```
    $ ./xmlsubscriber
    ```
**Note:** It will give list of options for running the sample. 

## Unit test and code coverage report

### Pre-requisite
Built ezmq-plus package.

### Run the unit test cases
1. Goto:  `~/protocol-ezmq-plus-go/`
2. Run the script:

   ```
   $ ./unittests.sh     : Native unit tests build for x86_64/armhf
   ```

### Generate the code coverage report
1. Goto `~/protocol-ezmq-plus-go/` </br>
2. Run the script:

   ```
   $ ./unittests.sh     : Native unit tests build for x86_64/armhf
   ```
3. Goto `~/protocol-ezmq-plus-go/src/go/ezmqx_unittests` </br>
4. Run the below command to open coverage report in web browser: </br>
     `$ go tool cover -html=coverage.out`

## Usage guide for ezmq-plus library (For micro-services) ##
1. The microservice which wants to use ezmq-plus Go library has to import ezmq package:
    `import go/ezmqx`
2. Reference ezmq-plus library APIs : [doc/godoc/ezmqx.html](doc/godoc/ezmqx.html)
3. Topic naming convention guide : [Naming Guide](https://github.sec.samsung.net/RS7-EdgeComputing/protocol-ezmq-plus-cpp/blob/1.0_rel/TOPIC_NAMING_CONVENTION.md)
