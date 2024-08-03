# TCC diagrams

```
@startuml
left to right direction
:writer: --> (Login)
:writer: --> (Criar Story)
:writer: --> (Criar Encounter)
:writer: --> (Criar Task)
:writer: --> (Vincula user como Storyteller)
:writer: --> (Vincula user como Player)
:writer: --> (Vincula Story em canal do Chat)
:writer: --> (Associa Task em Encounters)
:writer: --> (Associa Encounters com Encounters)
@enduml
```

```
@startuml
left to right direction
"Adiciona como user" as (join)
:user storyteller: --> (join)
:user storyteller: --> (Controla Encounters)
:user storyteller: --> (Controla NPCs)
:user player:  --> (join)
:user player:  --> (Controla Player)
@enduml

```

```
@startuml

' hide the spot
' hide circle

' avoid problems with angled crows feet
skinparam linetype ortho

entity "writer" as e01 {
  *id : number <<generated>>
}

entity "story" as e02 {
  *id : number <<generated>>
  ---
  *writer_id : number <<FK>>
}

entity "story_keys" as e03 {
  *id : number <<generated>>
  --
  *story_id : number <<FK>>
}

entity "access_story_keys" as e04 {
  *id : number <<generated>>
  --
  *story_keys_id : number <<FK>>
  *writer_id : number <<FK>>
}

entity "encounters" as e05 {
  *id : number <<generated>>
  --
  *story_id : number <<FK>>
  *writer_id : number <<FK>>
}

entity "slack_information" as e06 {
  *id : number <<generated>>
  --
}

entity "tasks" as e07 {
  *id : number <<generated>>
  --
}

entity "users" as e08 {
  *id : number <<generated>>
  --
}

entity "player" as e09 {
  *id : number <<generated>>
  --
  *stage_id : number <<FK>>
  *users_id : number <<FK>>
}

entity "stage" as e10 {
  *id : number <<generated>>
  --
  *story_id : number <<FK>>
  *storyteller_id : number <<FK>>
}

entity "stage_encounters" as e11 {
  *id : number <<generated>>
  --
  *stage_id : number <<FK>>
  *encounters_id : number <<FK>>
  *storyteller_id : number <<FK>>
}

entity "stage_running_tasks" as e12 {
  *id : number <<generated>>
  --
  *stage_id : number <<FK>>
  *stage_encounters_id : number <<FK>>
  *storyteller_id : number <<FK>>
  *tasks_id : number <<FK>>
}

entity "stage_next_encounter" as e13 {
  *id : number <<generated>>
  --
  *stage_id : number <<FK>>
  *current_encounter_id : number <<FK>>
  *next_encounter_id : number <<FK>>
}

entity "stage_encounter_activities" as e14 {
  *id : number <<generated>>
  --
  *stage_id : number <<FK>>
  *encounters_id : number <<FK>>
}

entity "non_players" as e15 {
  *id : number <<generated>>
  --
  *stage_id : number <<FK>>
  *storyteller_id : number <<FK>>
}

entity "stage_encounters_participants_players" as e16 {
  *id : number <<generated>>
  --
  *player_id : number <<FK>>
  *stage_encounters_id : number <<FK>>
}

entity "stage_encounters_participants_non_players" as e17 {
  *id : number <<generated>>
  --
  *storyteller_id : number <<FK>>
  *stage_encounters_id : number <<FK>>
}

entity "initiative" as e18 {
  *id : number <<generated>>
  --
  *stage_encounters_id : number <<FK>>
}

entity "initiative_participants" as e19 {
  *id : number <<generated>>
  --
  *initiative_id : number <<FK>>
}

entity "stage_channel" as e20 {
  *id : number <<generated>>
  --
  *stage_id : number <<FK>>
}

e01 ||..|{ e02
e01 ||..|{ e04
e01 ||..|{ e05

e03 ||..|| e04
e03 ||..|| e02

e02 ||..|{ e05

e08 ||..|| e10
e08 ||..|{ e09
e09 ||..|| e10

e10 ||..|{ e02
e10 ||..|{ e11
e11 ||..|| e05
e11 ||..|| e08

e10 ||..|{ e12
e12 ||..|| e11
e12 ||..|| e08
e12 ||..|| e07

e13 ||..|| e10

e14 ||..|| e10
e14 }|..|| e11

e15 ||..|| e10
e15 ||..|| e08

e16 ||..|{ e09
e16 }|..|| e11

e17 ||..|{ e15
e17 }|..|| e11

e18 ||..|| e11
e18 ||..|{ e19

e10 ||..|| e20

@enduml
```