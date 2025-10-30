import http from "k6/http";
import { check } from "k6";

function randomString(length) {
    const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    let result = "";
    for (let i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}

export const options = {
  scenarios: {
    rps_test: {
      executor:        "constant-arrival-rate",
      rate:            300,
      timeUnit:        "1s",
      duration:        "30s",
      preAllocatedVUs: 1000,
      maxVUs:          2000,
    },
  },
  thresholds: {
    "http_req_duration{name:shorten}": ["p(95)<600"],  // Realistic for high load
    "http_req_duration{name:resolve}": ["p(95)<150"],  // Cache hits should be fast
    "http_req_failed": ["rate<0.01"]  // Allow <1% failures
  },
};

export default function () {
  const BATCH_SIZE = 10; // Each VU sends 10 shorten + 10 resolve = 20 requests per iteration

  // Prepare batch of shorten requests
  const shortenRequests = [];
  const longUrls = [];

  for (let i = 0; i < BATCH_SIZE; i++) {
    const longUrl = `https://example.com/${randomString(32)}`;
    longUrls.push(longUrl);

    shortenRequests.push({
      method: "POST",
      url: "https://urlshortback.up.railway.app/api/shorten",
      body: JSON.stringify({ url: longUrl }),
      params: {
        headers: { "Content-Type": "application/json" },
        tags: { name: "shorten" }
      }
    });
  }

  // Batch shorten requests
  const shortenResponses = http.batch(shortenRequests);
  const codes = [];

  shortenResponses.forEach((res) => {
    check(res, {
      "status is 200": (res) => res.status === 200,
      "has code field": (res) => !!res.json("code"),
    });
    if (res.status === 200) {
      codes.push(res.json("code"));
    }
  });

  // Batch resolve requests
  const resolveRequests = codes.map((code) => ({
    method: "GET",
    url: `https://urlshortback.up.railway.app/${code}`,
    params: {
      tags: { name: "resolve" },
      redirects: 0
    }
  }));

  const resolveResponses = http.batch(resolveRequests);

  resolveResponses.forEach((res, idx) => {
    check(res, {
      "resolve status is 302": (res) => res.status === 302,
      "Location header is correct": (r) => {
        const loc = r.headers["Location"] || r.headers["location"];
        return loc === longUrls[idx];
      },
    });
  });
}
