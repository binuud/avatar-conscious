# Avatar State Diagram

## High Level Avatar State Transition

This is not finalized as yet.

```golang
enum AvatarStatus {
  AvatarActionStatusIgnore = 0;
  STATUS_DEFAULT = 1;
  // initialized and ready to accept commands
  STATUS_INIT = 2;

  STATUS_ACTIVE = 3;
  // waiting for user input
  // or some io

  STATUS_RERUN = 4;
    // Avatar is waiting for some input or IO
  // or an event
  STATUS_WAITING =5;
  STATUS_STOPPED = 12;
  STATUS_COMPLETED = 12;
}

enum AvatarResult {
  AvatarActionStatusIgnore = 0;
  RESULT_SUCCESS = 1;
  RESULT_FAILED = 2;
  RESULT_ERROR = 3;
}
```

```mermaid
stateDiagram-v2
    [*] --> STATUS_INIT : Start
    STATUS_INIT --> STATUS_ACTIVE
    STATUS_ACTIVE --> STATUS_WAITING : OnWaiting
    STATUS_ACTIVE --> STATUS_COMPLETED: OnComplete
    STATUS_ACTIVE --> STATUS_STOPPED : OnStop
    STATUS_WAITING --> STATUS_ACTIVE : OnResume
    STATUS_COMPLETED --> RESULT_SUCCESS: OnSuccess
    STATUS_COMPLETED --> RESULT_FAILED: OnFail
    STATUS_COMPLETED --> RESULT_ERROR: OnError
    STATUS_STOPPED --> RESULT_ERROR: OnError
```

```mermaid
stateDiagram-v2
    [*] --> Idle
    Idle --> Processing : Submit Request / Generate Invoice
    Processing --> Success : Payment Cleared / Dispatch Order
    Processing --> Failed : Timeout / Log Error
    Success --> [*]
    Failed --> [*]

```