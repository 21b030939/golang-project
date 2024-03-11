# Schedule project


## Schedule REST API
```
POST /shedules
GET /shedules/:id
PUT /shedules/:id
DELETE /shedules/:id
```

## DB Structure

```

Table schedule {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  discipline bigserial
  cabinet text
  time_period text
}

Table discipline {
  id bigserial [primary key]
  created_at timestamp
  updated_at timestamp
  name text
  description text
  credits text
}

Ref: schedules.discipline < discipline.id

```