# Protocol Documentation

<a name="top"></a>

## Table of Contents

- [app/sdk/proto/mlog/logs.proto](#app_sdk_proto_mlog_logs-proto)

  - [FileResponse](#logs-FileResponse)
  - [Log](#logs-Log)
  - [Log.MetadataEntry](#logs-Log-MetadataEntry)
  - [LogResponse](#logs-LogResponse)
  - [Logs](#logs-Logs)
  - [NewLog](#logs-NewLog)
  - [NewLog.MetadataEntry](#logs-NewLog-MetadataEntry)
  - [SearchQuery](#logs-SearchQuery)

  - [LogReader](#logs-LogReader)
  - [LogWriter](#logs-LogWriter)

- [Scalar Value Types](#scalar-value-types)

<a name="app_sdk_proto_mlog_logs-proto"></a>

<p align="right"><a href="#top">Top</a></p>

## app/sdk/proto/mlog/logs.proto

<a name="logs-FileResponse"></a>

### FileResponse

Resposta quando os logs são retornados como arquivo

| Field       | Type              | Label | Description                  |
| ----------- | ----------------- | ----- | ---------------------------- |
| file_url    | [string](#string) |       |                              |
| file_size   | [int64](#int64)   |       |                              |
| compression | [string](#string) |       | Tipo de compressão utilizada |

<a name="logs-Log"></a>

### Log

Mensagem com os logs retornados

| Field     | Type                                         | Label    | Description |
| --------- | -------------------------------------------- | -------- | ----------- |
| id        | [string](#string)                            |          |             |
| message   | [string](#string)                            |          |             |
| level     | [string](#string)                            |          |             |
| timestamp | [int64](#int64)                              |          |             |
| metadata  | [Log.MetadataEntry](#logs-Log-MetadataEntry) | repeated |             |

<a name="logs-Log-MetadataEntry"></a>

### Log.MetadataEntry

| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| key   | [string](#string) |       |             |
| value | [string](#string) |       |             |

<a name="logs-LogResponse"></a>

### LogResponse

Resposta ao registrar um log

| Field  | Type              | Label | Description |
| ------ | ----------------- | ----- | ----------- |
| id     | [string](#string) |       |             |
| status | [string](#string) |       |             |

<a name="logs-Logs"></a>

### Logs

Coleção de logs

| Field    | Type             | Label    | Description |
| -------- | ---------------- | -------- | ----------- |
| logs     | [Log](#logs-Log) | repeated |             |
| total    | [int32](#int32)  |          |             |
| has_more | [bool](#bool)    |          |             |

<a name="logs-NewLog"></a>

### NewLog

Mensagem para registrar um novo log

| Field     | Type                                               | Label    | Description |
| --------- | -------------------------------------------------- | -------- | ----------- |
| message   | [string](#string)                                  |          |             |
| level     | [string](#string)                                  |          |             |
| timestamp | [int64](#int64)                                    |          |             |
| metadata  | [NewLog.MetadataEntry](#logs-NewLog-MetadataEntry) | repeated |             |

<a name="logs-NewLog-MetadataEntry"></a>

### NewLog.MetadataEntry

| Field | Type              | Label | Description |
| ----- | ----------------- | ----- | ----------- |
| key   | [string](#string) |       |             |
| value | [string](#string) |       |             |

<a name="logs-SearchQuery"></a>

### SearchQuery

Consulta para buscar logs

| Field      | Type              | Label | Description                                      |
| ---------- | ----------------- | ----- | ------------------------------------------------ |
| start_time | [int64](#int64)   |       |                                                  |
| end_time   | [int64](#int64)   |       |                                                  |
| level      | [string](#string) |       |                                                  |
| page_size  | [int32](#int32)   |       |                                                  |
| page       | [int32](#int32)   |       |                                                  |
| as_file    | [bool](#bool)     |       | Se true, retorna como arquivo ao invés de stream |

<a name="logs-LogReader"></a>

### LogReader

Serviço para buscar logs

| Method Name  | Request Type                     | Response Type                      | Description                                                         |
| ------------ | -------------------------------- | ---------------------------------- | ------------------------------------------------------------------- |
| Search       | [SearchQuery](#logs-SearchQuery) | [Logs](#logs-Logs)                 | Busca logs retornando como stream de registros                      |
| ExportToFile | [SearchQuery](#logs-SearchQuery) | [FileResponse](#logs-FileResponse) | Busca logs retornando como arquivo                                  |
| StreamFile   | [SearchQuery](#logs-SearchQuery) | [Logs](#logs-Logs) stream          | Busca logs retornando como stream de chunks (para arquivos grandes) |

<a name="logs-LogWriter"></a>

### LogWriter

Serviço para registrar logs

| Method Name | Request Type           | Response Type                    | Description |
| ----------- | ---------------------- | -------------------------------- | ----------- |
| Register    | [NewLog](#logs-NewLog) | [LogResponse](#logs-LogResponse) |             |

## Scalar Value Types

| .proto Type                    | Notes                                                                                                                                           | C++    | Java       | Python      | Go      | C#         | PHP            | Ruby                           |
| ------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------- | ------ | ---------- | ----------- | ------- | ---------- | -------------- | ------------------------------ |
| <a name="double" /> double     |                                                                                                                                                 | double | double     | float       | float64 | double     | float          | Float                          |
| <a name="float" /> float       |                                                                                                                                                 | float  | float      | float       | float32 | float      | float          | Float                          |
| <a name="int32" /> int32       | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| <a name="int64" /> int64       | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| <a name="uint32" /> uint32     | Uses variable-length encoding.                                                                                                                  | uint32 | int        | int/long    | uint32  | uint       | integer        | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64     | Uses variable-length encoding.                                                                                                                  | uint64 | long       | int/long    | uint64  | ulong      | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32     | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s.                            | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64     | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s.                            | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| <a name="fixed32" /> fixed32   | Always four bytes. More efficient than uint32 if values are often greater than 2^28.                                                            | uint32 | int        | int         | uint32  | uint       | integer        | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64   | Always eight bytes. More efficient than uint64 if values are often greater than 2^56.                                                           | uint64 | long       | int/long    | uint64  | ulong      | integer/string | Bignum                         |
| <a name="sfixed32" /> sfixed32 | Always four bytes.                                                                                                                              | int32  | int        | int         | int32   | int        | integer        | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes.                                                                                                                             | int64  | long       | int/long    | int64   | long       | integer/string | Bignum                         |
| <a name="bool" /> bool         |                                                                                                                                                 | bool   | boolean    | boolean     | bool    | bool       | boolean        | TrueClass/FalseClass           |
| <a name="string" /> string     | A string must always contain UTF-8 encoded or 7-bit ASCII text.                                                                                 | string | String     | str/unicode | string  | string     | string         | String (UTF-8)                 |
| <a name="bytes" /> bytes       | May contain any arbitrary sequence of bytes.                                                                                                    | string | ByteString | str         | []byte  | ByteString | string         | String (ASCII-8BIT)            |
