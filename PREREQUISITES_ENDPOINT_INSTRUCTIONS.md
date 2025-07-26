# Instruções para adicionar o endpoint PrerequisitesByCourse

## 1. Adicione no arquivo `/proto/academictoken/subject/query.proto`

### No service Query, adicione este RPC (após SubjectsByCourse):

```protobuf
  // PrerequisitesByCourse lists all prerequisites for subjects in a specific course
  rpc PrerequisitesByCourse(QueryPrerequisitesByCourseRequest) returns (QueryPrerequisitesByCourseResponse) {
    option (google.api.http).get = "/academictoken/subject/prerequisites/course/{course_id}";
  }
```

### No final do arquivo, adicione estas mensagens:

```protobuf
// QueryPrerequisitesByCourseRequest is the request type for the Query/PrerequisitesByCourse RPC method
message QueryPrerequisitesByCourseRequest {
  string course_id = 1;
}

// QueryPrerequisitesByCourseResponse is the response type for the Query/PrerequisitesByCourse RPC method
message QueryPrerequisitesByCourseResponse {
  // Map of subject_id to its prerequisite groups
  map<string, PrerequisiteGroups> prerequisites = 1;
}

// PrerequisiteGroups wraps repeated PrerequisiteGroup for map usage
message PrerequisiteGroups {
  repeated PrerequisiteGroup groups = 1;
}
```

## 2. O handler já está criado em:
`/x/subject/keeper/query_prerequisites_by_course.go`

## 3. Após adicionar ao proto, execute:

```bash
cd /Users/biancamsp/Desktop/Academic_Token/academictoken/academictoken
make proto-gen
```

## 4. Reinicie o backend para que o novo endpoint esteja disponível

## 5. O endpoint estará disponível em:
`GET /academictoken/subject/prerequisites/course/{course_id}`

## Exemplo de resposta esperada:

```json
{
  "prerequisites": {
    "subject-1": {
      "groups": [
        {
          "id": "prereq-group-1",
          "subjectId": "subject-1",
          "groupType": "ALL",
          "subjectIds": ["subject-2", "subject-3"]
        }
      ]
    },
    "subject-4": {
      "groups": [
        {
          "id": "prereq-group-2",
          "subjectId": "subject-4",
          "groupType": "ANY",
          "subjectIds": ["subject-1", "subject-5"]
        }
      ]
    }
  }
}
```
