# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: protos/models/v1alpha1/common.proto
# Protobuf Python Version: 5.27.3
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    27,
    3,
    '',
    'protos/models/v1alpha1/common.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from protos.models.v1alpha1 import enums_pb2 as protos_dot_models_dot_v1alpha1_dot_enums__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n#protos/models/v1alpha1/common.proto\x12\x0fmodels.v1alpha1\x1a\"protos/models/v1alpha1/enums.proto\"M\n\x0bSearchParam\x12\x0c\n\x01q\x18\x01 \x01(\tR\x01q\x12\x30\n\x06search\x18\x02 \x03(\x0b\x32\x18.models.v1alpha1.MapListR\x06search\"0\n\x06Mapper\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n\x05value\x18\x02 \x01(\tR\x05value\"3\n\x07MapList\x12\x10\n\x03key\x18\x01 \x01(\tR\x03key\x12\x16\n\x06values\x18\x02 \x03(\tR\x06values\"X\n\x0e\x42\x61seIdentifier\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12\x12\n\x04type\x18\x02 \x01(\tR\x04type\x12\x1e\n\nidentifier\x18\x03 \x01(\tR\nidentifier\"t\n\x0cSyncSchedule\x12\x38\n\tfrequency\x18\x01 \x01(\x0e\x32\x1a.models.v1alpha1.FrequencyR\tfrequency\x12\x14\n\x05value\x18\x02 \x01(\x03R\x05value\x12\x14\n\x05limit\x18\x03 \x01(\x03R\x05limit\"\\\n\x0bScheduleTab\x12\x37\n\x07syncTab\x18\x01 \x03(\x0b\x32\x1d.models.v1alpha1.SyncScheduleR\x07syncTab\x12\x14\n\x05limit\x18\x02 \x01(\x03R\x05limit\"\xe1\x02\n\tJWTParams\x12\x12\n\x04name\x18\x01 \x01(\tR\x04name\x12\"\n\x0cpublicJWTKey\x18\x02 \x01(\tR\x0cpublicJWTKey\x12$\n\rprivateJWTKey\x18\x03 \x01(\tR\rprivateJWTKey\x12\x10\n\x03vId\x18\x04 \x01(\tR\x03vId\x12K\n\x10signingAlgorithm\x18\x05 \x01(\x0e\x32\x1f.models.v1alpha1.EncryptionAlgoR\x10signingAlgorithm\x12\x30\n\x13isAlreadyInSecretBS\x18\x06 \x01(\x08R\x13isAlreadyInSecretBS\x12\x35\n\x06status\x18\x07 \x01(\x0e\x32\x1d.models.v1alpha1.CommonStatusR\x06status\x12.\n\x12generateInPlatform\x18\x08 \x01(\x08R\x12generateInPlatform\"\xad\x01\n\tTlsConfig\x12\x32\n\x07tlsType\x18\x01 \x01(\x0e\x32\x18.models.v1alpha1.TlsTypeR\x07tlsType\x12\x1e\n\ncaCertFile\x18\x02 \x01(\tR\ncaCertFile\x12$\n\rserverKeyFile\x18\x03 \x01(\tR\rserverKeyFile\x12&\n\x0eserverCertFile\x18\x04 \x01(\tR\x0eserverCertFile\"\xb7\x01\n\x08SSLCerts\x12\x18\n\x07sslCert\x18\x01 \x01(\tR\x07sslCert\x12\x16\n\x06sslKey\x18\x02 \x01(\tR\x06sslKey\x12\x10\n\x03vId\x18\x03 \x01(\tR\x03vId\x12\x30\n\x13isAlreadyInSecretBS\x18\x04 \x01(\x08R\x13isAlreadyInSecretBS\x12\x35\n\x06status\x18\x05 \x01(\x0e\x32\x1d.models.v1alpha1.CommonStatusR\x06status\"\xbf\x01\n\nDMResponse\x12\x18\n\x07message\x18\x01 \x01(\tR\x07message\x12\"\n\x0c\x64mStatusCode\x18\x02 \x01(\tR\x0c\x64mStatusCode\x12>\n\rcustomMessage\x18\x03 \x03(\x0b\x32\x18.models.v1alpha1.MapListR\rcustomMessage\x12\x33\n\x06\x61gents\x18\x04 \x03(\x0b\x32\x1b.models.v1alpha1.AgentShortR\x06\x61gents\"\x8e\x01\n\nAgentShort\x12\x14\n\x05\x61gent\x18\x01 \x01(\tR\x05\x61gent\x12\x1c\n\tagentType\x18\x02 \x01(\tR\tagentType\x12\x1e\n\ngoalStatus\x18\x03 \x01(\tR\ngoalStatus\x12\x14\n\x05\x65rror\x18\x04 \x03(\tR\x05\x65rror\x12\x16\n\x06result\x18\x05 \x01(\tR\x06result\"\x10\n\x0e\x44ynamicMessage\"K\n\x14\x44ynamicMessageUpdate\x12\x33\n\x04\x62ody\x18\x01 \x01(\x0b\x32\x1f.models.v1alpha1.DynamicMessageR\x04\x62ody\"\xd4\x02\n\x10VapusCredentials\x12\x1a\n\x08username\x18\x01 \x01(\tR\x08username\x12\x1a\n\x08password\x18\x02 \x01(\tR\x08password\x12\x1a\n\x08\x61piToken\x18\x03 \x01(\tR\x08\x61piToken\x12\x41\n\x0c\x61piTokenType\x18\x04 \x01(\x0e\x32\x1d.models.v1alpha1.ApiTokenTypeR\x0c\x61piTokenType\x12\x35\n\x08\x61wsCreds\x18\x05 \x01(\x0b\x32\x19.models.v1alpha1.AWSCredsR\x08\x61wsCreds\x12\x35\n\x08gcpCreds\x18\x06 \x01(\x0b\x32\x19.models.v1alpha1.GCPCredsR\x08gcpCreds\x12;\n\nazureCreds\x18\x07 \x01(\x0b\x32\x1b.models.v1alpha1.AzureCredsR\nazureCreds\"\x92\x01\n\x08\x41WSCreds\x12 \n\x0b\x61\x63\x63\x65ssKeyId\x18\x01 \x01(\tR\x0b\x61\x63\x63\x65ssKeyId\x12(\n\x0fsecretAccessKey\x18\x02 \x01(\tR\x0fsecretAccessKey\x12\x16\n\x06region\x18\x03 \x01(\tR\x06region\x12\"\n\x0csessionToken\x18\x04 \x01(\tR\x0csessionToken\"\xa8\x01\n\x08GCPCreds\x12,\n\x11serviceAccountKey\x18\x01 \x01(\tR\x11serviceAccountKey\x12$\n\rbase64Encoded\x18\x02 \x01(\x08R\rbase64Encoded\x12\x16\n\x06region\x18\x03 \x01(\tR\x06region\x12\x1c\n\tprojectId\x18\x04 \x01(\tR\tprojectId\x12\x12\n\x04zone\x18\x05 \x01(\tR\x04zone\"h\n\nAzureCreds\x12\x1a\n\x08tenantId\x18\x01 \x01(\tR\x08tenantId\x12\x1a\n\x08\x63lientId\x18\x02 \x01(\tR\x08\x63lientId\x12\"\n\x0c\x63lientSecret\x18\x03 \x01(\tR\x0c\x63lientSecret\"S\n\tDigestVal\x12.\n\x04\x61lgo\x18\x01 \x01(\x0e\x32\x1a.models.v1alpha1.HashAlgosR\x04\x61lgo\x12\x16\n\x06\x64igest\x18\x02 \x01(\tR\x06\x64igest\"\x0e\n\x0c\x45mptyRequest\"\xe2\x01\n\tAuthnOIDC\x12\x1a\n\x08\x63\x61llback\x18\x01 \x01(\tR\x08\x63\x61llback\x12\x1a\n\x08\x63lientId\x18\x02 \x01(\tR\x08\x63lientId\x12\"\n\x0c\x63lientSecret\x18\x03 \x01(\tR\x0c\x63lientSecret\x12\x10\n\x03vId\x18\x04 \x01(\tR\x03vId\x12\x30\n\x13isAlreadyInSecretBS\x18\x05 \x01(\x08R\x13isAlreadyInSecretBS\x12\x35\n\x06status\x18\x06 \x01(\x0e\x32\x1d.models.v1alpha1.CommonStatusR\x06status\"d\n\x04Logs\x12\x12\n\x04time\x18\x01 \x01(\x03R\x04time\x12\x14\n\x05hTime\x18\x02 \x01(\tR\x05hTime\x12\x18\n\x07logType\x18\x03 \x01(\tR\x07logType\x12\x18\n\x07message\x18\x04 \x01(\tR\x07messageB\xca\x01\n\x13\x63om.models.v1alpha1B\x0b\x43ommonProtoP\x01ZIgithub.com/vapusdata-ecosystem/vapusai-studio/apis/protos/models/v1alpha1\xa2\x02\x03MXX\xaa\x02\x0fModels.V1alpha1\xca\x02\x0fModels\\V1alpha1\xe2\x02\x1bModels\\V1alpha1\\GPBMetadata\xea\x02\x10Models::V1alpha1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'protos.models.v1alpha1.common_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\023com.models.v1alpha1B\013CommonProtoP\001ZIgithub.com/vapusdata-ecosystem/vapusai-studio/apis/protos/models/v1alpha1\242\002\003MXX\252\002\017Models.V1alpha1\312\002\017Models\\V1alpha1\342\002\033Models\\V1alpha1\\GPBMetadata\352\002\020Models::V1alpha1'
  _globals['_SEARCHPARAM']._serialized_start=92
  _globals['_SEARCHPARAM']._serialized_end=169
  _globals['_MAPPER']._serialized_start=171
  _globals['_MAPPER']._serialized_end=219
  _globals['_MAPLIST']._serialized_start=221
  _globals['_MAPLIST']._serialized_end=272
  _globals['_BASEIDENTIFIER']._serialized_start=274
  _globals['_BASEIDENTIFIER']._serialized_end=362
  _globals['_SYNCSCHEDULE']._serialized_start=364
  _globals['_SYNCSCHEDULE']._serialized_end=480
  _globals['_SCHEDULETAB']._serialized_start=482
  _globals['_SCHEDULETAB']._serialized_end=574
  _globals['_JWTPARAMS']._serialized_start=577
  _globals['_JWTPARAMS']._serialized_end=930
  _globals['_TLSCONFIG']._serialized_start=933
  _globals['_TLSCONFIG']._serialized_end=1106
  _globals['_SSLCERTS']._serialized_start=1109
  _globals['_SSLCERTS']._serialized_end=1292
  _globals['_DMRESPONSE']._serialized_start=1295
  _globals['_DMRESPONSE']._serialized_end=1486
  _globals['_AGENTSHORT']._serialized_start=1489
  _globals['_AGENTSHORT']._serialized_end=1631
  _globals['_DYNAMICMESSAGE']._serialized_start=1633
  _globals['_DYNAMICMESSAGE']._serialized_end=1649
  _globals['_DYNAMICMESSAGEUPDATE']._serialized_start=1651
  _globals['_DYNAMICMESSAGEUPDATE']._serialized_end=1726
  _globals['_VAPUSCREDENTIALS']._serialized_start=1729
  _globals['_VAPUSCREDENTIALS']._serialized_end=2069
  _globals['_AWSCREDS']._serialized_start=2072
  _globals['_AWSCREDS']._serialized_end=2218
  _globals['_GCPCREDS']._serialized_start=2221
  _globals['_GCPCREDS']._serialized_end=2389
  _globals['_AZURECREDS']._serialized_start=2391
  _globals['_AZURECREDS']._serialized_end=2495
  _globals['_DIGESTVAL']._serialized_start=2497
  _globals['_DIGESTVAL']._serialized_end=2580
  _globals['_EMPTYREQUEST']._serialized_start=2582
  _globals['_EMPTYREQUEST']._serialized_end=2596
  _globals['_AUTHNOIDC']._serialized_start=2599
  _globals['_AUTHNOIDC']._serialized_end=2825
  _globals['_LOGS']._serialized_start=2827
  _globals['_LOGS']._serialized_end=2927
# @@protoc_insertion_point(module_scope)
