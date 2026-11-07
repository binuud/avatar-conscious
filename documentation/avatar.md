# Avatar

## Avatar Store Files

* avatar-name.md md file containing front matter and optional steps as MD document
* avatar-name-steps.yaml containing avatar steps in yaml format, if MD document is given this file is not required, this save compilation time, since LLM is not used to determine the steps, and is much faster avatar runs
* avatar-name-ui.yaml collection of ui elements used by avatar, later to be created by LLM itself
* avatar-name.png

## avatar-name-steps.yaml file

This file supports slash commands
eg: /launch-ui [ui-name] to launch ui directly, does not use LLM to determine which UI to launch, resutls in faster user experience


Step Types
```golang
  enum StepType {
    STEP_TYPE_DEFAULT = 0;
    STEP_TYPE_AI = 1;
    STEP_TYPE_UI = 2;
  }
```