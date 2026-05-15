import http from 'k6/http';
import { check } from 'k6';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:7080';
const PRODUCT_ID = __ENV.PRODUCT_ID || 'PROD-001';

export const options = {
  scenarios: {
    warm_up: {
      executor: 'constant-arrival-rate',
      rate: 100,
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 100,
      maxVUs: 300,
      exec: 'reviewSummary',
    },
    review_summary_1000_rps: {
      executor: 'constant-arrival-rate',
      rate: 1000,
      timeUnit: '1s',
      duration: '30s',
      startTime: '30s',
      preAllocatedVUs: 600,
      maxVUs: 2000,
      exec: 'reviewSummary',
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<2000'],
  },
};

export function reviewSummary() {
  const response = http.get(`${BASE_URL}/api/v1/products/${PRODUCT_ID}/review-summary`);

  check(response, {
    'status is 200': (res) => res.status === 200,
  });
}
