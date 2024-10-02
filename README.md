Data Generator with Gemini (Cloud Run functions)
-----------------------------
Details TBA

# TODO
1. //TODO: add function calling to parse the first response (e.g., to set max # of rows to be generated) and then process the actual request
   * See examples/main.go for details
```
Received request:
"Create 149 dummy cookie recipes using this JSON schema: Recipe = {'recipeName': string} Return: Array<Recipe>"

Received function call response:
{"promptParser" map["counter":%!q(float64=8) "prompt":"Create 20 dummy cookie recipes using this JSON schema: Recipe = {\\'recipeName\\': string} Return: Array<Recipe>"]}

Calling actual function for one time only:
[{"recipeName": "Chocolate Chip Cookies"}, {"recipeName": "Oatmeal Raisin Cookies"}, {"recipeName": "Peanut Butter Cookies"}, {"recipeName": "Snickerdoodles"}, {"recipeName": "Sugar Cookies"}, {"recipeName": "Gingerbread Cookies"}, {"recipeName": "Shortbread Cookies"}, {"recipeName": "Macarons"}, {"recipeName": "Brownies"}, {"recipeName": "Biscotti"}, {"recipeName": "Lemon Bars"}, {"recipeName": "Toffee Cookies"}, {"recipeName": "Gingersnaps"}, {"recipeName": "Peanut Butter Blossoms"}, {"recipeName": "Thumbprint Cookies"}, {"recipeName": "Linzer Cookies"}, {"recipeName": "Rugelach"}, {"recipeName": "Alfajores"}, {"recipeName": "Madeleines"}, {"recipeName": "Palmiers"}]

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
