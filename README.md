Data Generator with Gemini (Cloud Run functions)
-----------------------------
Details TBA

# TODO
1. //TODO: add function calling to parse the first response (e.g., to set max # of rows to be generated) and then process the actual request
   * See examples/main.go for details
```
Received request:
"Generate 100 rows of dummy data for a table that has id, name, and email columns."

Received function call response:
{"promptParser" map["prompt":"Generate a table with id, name, and email columns." "rows":%!q(float64=100)]}

Executing function call response:
map["prompt":"Generate a table with id, name, and email columns." "rows":%!q(float64=100)]

Received response:
{promptParser map[prompt:Generate a table with id, name and email columns rows:100]}

-- parallelize the data generation process --
```
2. //TODO: add SQLite/DuckDB integration to dedup the data

# How to run
## Run locally
```
FUNCTION_TARGET=DataGeneratorGemini PROJECT_ID=YOUR_PROJECT_ID LOCATION=YOUR_LOCATION CONTEXT_TIMEOUT_S=YOUR_CONTEXT_TIMEOUT_S go run cmd/main.go
```

## Example request (schema need to be included in the prompt)
```
curl -m 60 -X POST localhost:8080 \
-H "Content-Type: application/json" \
-d '{
  "requestId": "",
  "promptInput": "List a few popular cookie recipes using this JSON schema: Recipe = {'recipeName': string} Return: Array<Recipe>",
  "model": "gemini-1.5-flash-002",
  "modelConfig": {"temperature":0.2,"maxOutputTokens":8000,"topP":0.8,"topK":40}
}'
```

## Example response (enforced as JSON array and will throw error if not)
```
{
    "data": [
        {"recipeName": "Chocolate Chip Cookies"},
        {"recipeName": "Oatmeal Raisin Cookies"},
        {"recipeName": "Peanut Butter Cookies"},
        {"recipeName": "Sugar Cookies"}
    ],
    "errorMessage": ""
}
```
```
{"data":null,"errorMessage":"invalid response: not an array"}
```

# Additional notes
TBA
