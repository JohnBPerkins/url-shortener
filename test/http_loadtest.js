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
      executor: "constant-arrival-rate",
      rate: 2500,
      timeUnit: "1s",
      duration: "30s",
      preAllocatedVUs: 1000,
      maxVUs: 2000,
    },
  },
  thresholds: {
    "http_req_duration{name:shorten}": ["p(95)<50"],
    "http_req_duration{name:resolve}": ["p(95)<50"]
  },
};

export default function () {
  const longUrl = `https://example.com/${randomString(32)}`;
  const payload = JSON.stringify({ url: longUrl });

  const params = {
    headers: { "Content-Type": "application/json" },
    tags: {name: "shorten"}
  };

  let res = http.post("http://localhost:8080/api/shorten", payload, params);
  check(res, {
    "status is 200": (res) => res.status === 200,
    "has code field": (res) => !!res.json("code"),
  });

  const code = res.json("code");
  res = http.get(
    `http://localhost:8080/${code}`,
    {
      tags: {name: "resolve"},
      redirects: 0
    }
  );

  check(res, {
    "resolve status is 302": (res) => res.status === 302,
    "Location header is correct": (r) => {
      const loc = r.headers["Location"] || r.headers["location"];
      return loc === longUrl;
    },
  });
}
