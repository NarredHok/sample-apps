1. Health Check
Check if the service is running and healthy.
Endpoint: GET /health
curl http://localhost:8080/health
2. List All Patients
Get a list of all patients in the system.
Endpoint: GET /api/patients
curl http://localhost:8080/api/patients
3. Get Patient by Name
Retrieve information for a specific patient by their name.
Endpoint: GET /api/patients/{name}
Note: The name is case-insensitive. Spaces should be URL-encoded as %20.
curl http://localhost:8080/api/patients/Johnson%20Fake
4. Create New Patient
Add a new patient to the system.
Endpoint: POST /api/patients
Headers:
Content-Type: application/json

Required Fields:
name (string)
dateOfBirth (string, format: YYYY-MM-DD)
email (string)

Optional Fields:
gender (string)
illness (string)

curl -X POST http://localhost:8080/api/patients \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Alice Williams",
    "dateOfBirth": "1988-05-20",
    "gender": "Female",
    "illness": "Allergies",
    "email": "alice.williams@email.com"
  }'
