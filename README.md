# mcp-api-tester Tools & Interfaces

This document summarizes potential tools or interfaces that enable an LLM (Large Language Model) to perform automated API testing. By leveraging these tools, the LLM can read API documentation, select and test specific endpoints using `net/http`, and help developers quickly identify potential issues.

## Core Tools (Minimal Viable Tools)

The following three tools are the most basic and essential interfaces, allowing the LLM to conduct **minimum viable** automated testing:

1. **listAllAPIFromDocument**

   - **Purpose**: Lists all available APIs (name, URL, HTTP method, etc.) from the documentation.
   - **Example**: A sample return value might look like:
     ```json
     [
       { "name": "getUser", "method": "GET", "url": "/users/{id}" },
       { "name": "createUser", "method": "POST", "url": "/users" }
       ...
     ]
     ```
   - **When the LLM Uses It**: To browse available APIs and choose which to test.

2. **getSingleAPIDetail**

   - **Purpose**: Retrieves detailed documentation of a specific API by name or path, such as:
     - Request parameters (query, path, body, etc.)
     - Response structure (schema)
     - Possible status codes and error definitions
   - **Example**: A sample return value might look like:
     ```json
     {
       "name": "getUser",
       "method": "GET",
       "url": "/users/{id}",
       "parameters": [
         { "in": "path", "name": "id", "type": "string", "required": true }
       ],
       "responses": {
         "200": {
           "description": "Returns user data on success",
           "schema": { ... }
         },
         "404": {
           "description": "User not found"
         }
       }
     }
     ```
   - **When the LLM Uses It**: To generate test parameters and expected responses for a specific endpoint.

3. **callNetHTTP**
   - **Purpose**: Allows the LLM to send real HTTP requests and receive results (status code, headers, body).
   - **Example**:
     ```json
     {
       "request": {
         "method": "POST",
         "url": "/users",
         "headers": { "Content-Type": "application/json" },
         "body": { "name": "NewUser", "age": 30 }
       },
       "response": {
         "status": 201,
         "headers": { ... },
         "body": { "id": "123", "name": "NewUser", "age": 30 }
       }
     }
     ```
   - **When the LLM Uses It**: To execute test calls and compare actual results against expected outcomes.

---

## Advanced Suggested Tools (Enhanced Tools)

These tools are not mandatory, but they can significantly improve automated testing, reporting, environment setup, and more.

### 1. Response Validation / Parsing Tools

1. **Schema Validation**

   - `validateResponseWithSchema(responseBody, schemaInfo) -> (isValid, errorList)`
   - Compares the returned JSON body against the documented schema to check for field consistency, data types, and any required fields.

2. **JSON Parsing**
   - `parseJSON(responseBody) -> (parsedObject)`
   - Enables the LLM to conveniently read and inspect specific fields in the response (e.g., confirm that certain fields are not empty).

### 2. Context / Logging Tools for Test Execution

1. **Test Result Logging / Storage**

   - `storeTestResult(testCase, response, passOrFail, errorMessage)`
   - Lets the LLM record execution details, errors, and outcomes for each test case, making it easier to compile reports.

2. **Historical Test Results**
   - `getTestHistory(apiName or endpoint) -> (historyData)`
   - Allows the LLM (or other systems) to query past test results for the same endpoint and understand regressions or changes over time.

### 3. Environment Setup / Management Tools

1. **Authentication / Credential Setup**

   - `getAuthToken(username, password) -> token`
   - Some APIs require a token or cookie for authentication. The LLM may need to call this tool first to obtain proper credentials.

2. **Reset / Initialize Test Data**
   - `resetTestData()` or `seedTestData()`
   - Ensures the database or system state is in a known default condition before testing begins.

### 4. Fuzz / Edge Case Testing Tools

- **generateFuzzyInput(fieldType, constraints)**
- Returns extreme or unusual input (e.g., negative values, extremely large values, random characters, special symbols) for stress-testing and security checks.

---

## Recommended Development Flow

1. **Integrate the Minimal Viable Tools**

   1. Implement `listAllAPIFromDocument`, `getSingleAPIDetail`, and `callNetHTTP` in your system so the LLM can:
      - Obtain a list of available endpoints
      - Fetch detailed specs for a given API
      - Make real HTTP requests

2. **Connect the LLM**

   - Enable the LLM to use the above three tools. It can list potential endpoints, retrieve details for each, generate test scenarios, and finally run them via `callNetHTTP`.

3. **Add Validation Tools (Response / Schema)**

   - Provide utilities that let the LLM verify whether returned data matches the documented schema and expectations.

4. **Expand on Test Logging & Environment Management**

   - Offer interfaces for storing test results, querying historical tests, resetting the environment, and retrieving auth tokens.
   - The LLM can call these tools to maintain a consistent state and documentation of test processes.

5. **Fuzz / Edge Case Testing**
   - Integrate a tool to generate malicious or extreme values, or let the LLM generate them automatically to conduct in-depth stress or security testing.

---

## Important Notes

- **Security & Permissions**

  - When the LLM makes API calls, ensure these tests run in a safe environment (development or sandbox) and do not compromise sensitive data.

- **Accuracy of LLM-Generated Content**

  - The LLM may produce invalid or incorrect parameters. Additional validation or error handling may be needed in your interfaces.

- **Cannot Completely Replace Traditional Testing**
  - While the LLM can quickly produce boundary and abnormal test cases, traditional unit tests remain essential for thorough coverage.
