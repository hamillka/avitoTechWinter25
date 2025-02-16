import http from 'k6/http';
import {check, sleep} from 'k6';

export let options = {
    scenarios: {
        constant_request_rate: {
            executor: 'constant-arrival-rate',
            rate: 1000,
            timeUnit: '1s',
            duration: '5m',
            preAllocatedVUs: 75,
            maxVUs: 300,
            gracefulStop: '30s',
        },
    },
    thresholds: {
        http_req_duration: ['p(95)<50'],
        http_req_failed: ['rate<0.0001'],
    },
};

export default function () {
    let res = http.get('http://localhost:8080/api/info', {
        headers: {
            'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NzExNzEzNzksInVzZXJuYW1lIjoidXNlcjk5MjAwIn0.eBGxEUYZFpiJu1AuRnCd71eW4GTD0zW3_UDihIbBVhc',
        },
    });

    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(0.001);
}