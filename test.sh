access_token="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NUb2tlbiI6eyJkYXRhIjp7ImVtYWlsIjoicnpQRUVSUHRqTjZXN0Y0NE9iQnAzSDczUlhyb2ZqcHdRZzFjRWJkVHBIV1lVYnExdms0VnZaak5FTUttMFlyWmFQd05ZZnVUNllJQmZCUkZtOEk0bzBuOW44R3hRdU41cVVpS0BnbWFpbC5jb20iLCJfaWQiOiI2NWQxODliMDE1ZWVkODRhYjkyZWEyOWEifSwiZW5jcnlwdGVkX2RhdGEiOiJlUmxZZXhLWHd1SGl2RCt0VEwwPSIsInVuYW1lIjoicnpQRUVSUHRqTjZXN0Y0NE9iQnAzSDczUlhyb2ZqcHdRZzFjRWJkVHBIV1lVYnExdms0VnZaak5FTUttMFlyWmFQd05ZZnVUNllJQmZCUkZtOEk0bzBuOW44R3hRdU41cVVpS0BnbWFpbC5jb20iLCJ0b2tlbl9pZCI6IjY1ZDE4OWIwMTVlZWQ4NGFiOTJlYTI5YV9XbWFWMHVFakUrQXFHUUVBUTZ0a09wRlV1NnpXTDZHbFJmZ1pidktiN0UycUE1Z0JvV3hPNjlGYUlVU0hlcmZRaE41UGZJOUkzMGJtRTIvYlZuRXNVbEQwbTJROFVTQ1owUVFUdVA1ckZCY1RDamUxRngxaFo2MXNINEE1R1YyK2xvL2l0Zz09XzE3MDgyMzEwOTQxMTEiLCJleHAiOjE3MDgyMzQ2OTQxMTEsImlzc3VlZF9hdCI6MTcwODIzMTA5NDExMSwiY3NyZl90b2tlbiI6ImVDQk1hbCtkN3VQWHFSK0xWcmg0NGxvbGtDaFFBTEp1TWwvNkZ1MlJPdEp1VVlmdkN4cnJ2elVWd0hJSmR3NG1SY3VSZ1JBc0JWU3ByN09CWnFKMVQ1M1BmbjZ4bWFNby9VZUpWQnVTcEhPelFRQ1RiQ0N6UGU5dnNjS01SSER6cHo1SnlqS1kxODBMQThGZTFxcGYybWJTbUNBem5hKzlLbG1KSHRFakw3aVFTSGUxMk4zVHFRPT0ifX0.bJ2-1DSjr-P7KAUpwnsLLgE12JhGdj4o73AS4qvGJwQ"

csrf_token="rVytmyMdP5DGndtPyGvyU91faq78IpInOjGqKIXnJ8j/OeH8tAszsM3sRNzJgIMWfUOQgge3mASZBdJfbGBSGOXuuv5hFmNSwHuTY9A6X2G6rDfXvFhL+rAkvGe5K+W2Tk7H9A=="



curl  -iv --cookie  "access_token=$access_token" --header "Csrf_token: $csrf_token"  "http://localhost:8000/api/user/" 