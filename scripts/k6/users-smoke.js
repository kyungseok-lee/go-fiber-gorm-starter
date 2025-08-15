import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
  stages: [
    { duration: '10s', target: 10 },  // Ramp up to 10 users
    { duration: '30s', target: 50 },  // Stay at 50 users for 30 seconds
    { duration: '10s', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
    http_req_failed: ['rate<0.01'],   // Error rate must be less than 1%
    errors: ['rate<0.01'],            // Custom error rate must be less than 1%
  },
};

// Base URL - can be overridden with environment variable
const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

// Test data
const testUsers = [
  {
    name: 'Test User 1',
    email: `test.user.1.${Date.now()}@example.com`,
    status: 'active'
  },
  {
    name: 'Test User 2', 
    email: `test.user.2.${Date.now()}@example.com`,
    status: 'active'
  }
];

let createdUserIds = [];

export default function () {
  // Health check
  healthCheck();
  
  // User CRUD operations
  testCreateUser();
  testListUsers();
  testGetUser();
  testUpdateUser();
  
  // Add some think time between iterations
  sleep(1);
}

// Setup function - runs once at the beginning
export function setup() {
  console.log('Starting smoke test...');
  console.log(`Base URL: ${BASE_URL}`);
  
  // Verify API is accessible
  const healthResponse = http.get(`${BASE_URL}/health`);
  check(healthResponse, {
    'Health check passes': (r) => r.status === 200,
  });
  
  return { baseUrl: BASE_URL };
}

// Teardown function - runs once at the end
export function teardown(data) {
  console.log('Cleaning up test data...');
  
  // Clean up created users
  createdUserIds.forEach(userId => {
    const deleteResponse = http.del(`${BASE_URL}/v1/users/${userId}`);
    check(deleteResponse, {
      'User cleanup successful': (r) => r.status === 204,
    });
  });
  
  console.log('Smoke test completed.');
}

function healthCheck() {
  const response = http.get(`${BASE_URL}/health`);
  
  const success = check(response, {
    'Health status is 200': (r) => r.status === 200,
    'Health response time < 100ms': (r) => r.timings.duration < 100,
    'Health response contains status': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && body.data.status === 'ok';
      } catch (e) {
        return false;
      }
    },
  });
  
  errorRate.add(!success);
}

function testCreateUser() {
  const userData = testUsers[Math.floor(Math.random() * testUsers.length)];
  
  // Create unique email for this test run
  userData.email = `test.${__VU}.${Date.now()}@example.com`;
  
  const response = http.post(
    `${BASE_URL}/v1/users`,
    JSON.stringify(userData),
    {
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );
  
  const success = check(response, {
    'Create user status is 201': (r) => r.status === 201,
    'Create user response time < 500ms': (r) => r.timings.duration < 500,
    'Create user response contains user data': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && body.data.id && body.data.email === userData.email;
      } catch (e) {
        return false;
      }
    },
  });
  
  if (success && response.status === 201) {
    try {
      const body = JSON.parse(response.body);
      if (body.data && body.data.id) {
        createdUserIds.push(body.data.id);
      }
    } catch (e) {
      console.error('Failed to parse create user response:', e);
    }
  }
  
  errorRate.add(!success);
}

function testListUsers() {
  const response = http.get(`${BASE_URL}/v1/users?offset=0&limit=10`);
  
  const success = check(response, {
    'List users status is 200': (r) => r.status === 200,
    'List users response time < 300ms': (r) => r.timings.duration < 300,
    'List users response contains pagination': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && body.pagination && 
               typeof body.pagination.total === 'number';
      } catch (e) {
        return false;
      }
    },
  });
  
  errorRate.add(!success);
}

function testGetUser() {
  // Only test if we have created users
  if (createdUserIds.length === 0) return;
  
  const userId = createdUserIds[Math.floor(Math.random() * createdUserIds.length)];
  const response = http.get(`${BASE_URL}/v1/users/${userId}`);
  
  const success = check(response, {
    'Get user status is 200': (r) => r.status === 200,
    'Get user response time < 200ms': (r) => r.timings.duration < 200,
    'Get user response contains user data': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && body.data.id == userId;
      } catch (e) {
        return false;
      }
    },
  });
  
  errorRate.add(!success);
}

function testUpdateUser() {
  // Only test if we have created users
  if (createdUserIds.length === 0) return;
  
  const userId = createdUserIds[Math.floor(Math.random() * createdUserIds.length)];
  const updateData = {
    name: `Updated User ${Date.now()}`,
    status: 'inactive'
  };
  
  const response = http.put(
    `${BASE_URL}/v1/users/${userId}`,
    JSON.stringify(updateData),
    {
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );
  
  const success = check(response, {
    'Update user status is 200': (r) => r.status === 200,
    'Update user response time < 400ms': (r) => r.timings.duration < 400,
    'Update user response contains updated data': (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.data && body.data.name === updateData.name;
      } catch (e) {
        return false;
      }
    },
  });
  
  errorRate.add(!success);
}

// Test invalid scenarios
function testErrorScenarios() {
  // Test 404 for non-existent user
  const notFoundResponse = http.get(`${BASE_URL}/v1/users/999999`);
  check(notFoundResponse, {
    'Non-existent user returns 404': (r) => r.status === 404,
  });
  
  // Test 400 for invalid user creation
  const invalidUserResponse = http.post(
    `${BASE_URL}/v1/users`,
    JSON.stringify({ name: '' }), // Invalid: empty name
    {
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );
  check(invalidUserResponse, {
    'Invalid user creation returns 400': (r) => r.status === 400,
  });
}

// Additional test scenarios for stress testing
export function stressTest() {
  // Run multiple operations concurrently
  const promises = [];
  
  for (let i = 0; i < 5; i++) {
    promises.push(testCreateUser());
    promises.push(testListUsers());
  }
  
  // Wait for all operations to complete
  Promise.all(promises);
}

/*
Usage examples:

1. Basic smoke test:
   k6 run scripts/k6/users-smoke.js

2. With custom base URL:
   BASE_URL=http://your-api.com k6 run scripts/k6/users-smoke.js

3. With different user load:
   k6 run --vus 100 --duration 2m scripts/k6/users-smoke.js

4. Generate detailed report:
   k6 run --out json=results.json scripts/k6/users-smoke.js

5. Run with specific thresholds:
   k6 run --threshold http_req_duration=p(95)<1000 scripts/k6/users-smoke.js

Key metrics to monitor:
- http_req_duration: Request duration (check p95, p99)
- http_req_failed: Failed request rate
- http_reqs: Request rate (RPS)
- vus: Virtual users
- iteration_duration: Full iteration time

Performance benchmarks:
- P95 response time should be < 500ms for most endpoints
- Error rate should be < 1%
- API should handle at least 50 RPS without degradation
- Memory usage should remain stable during load
*/