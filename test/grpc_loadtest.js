import grpc from "k6/net/grpc";
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
      rate:            1000,
      timeUnit:        "1s",
      duration:        "30s",
      preAllocatedVUs: 1000,
      maxVUs:          2000,
    },
  },
  thresholds: {
    "grpc_req_duration{method:Shorten}": ["p(95)<100"],
    "grpc_req_duration{method:Resolve}": ["p(95)<100"],
  },
};

const client = new grpc.Client();
client.load(["../proto"], "shortener.proto");

let connected = false;

export default function () {
  if (!connected) {
    client.connect("localhost:50051", { plaintext: true });
    connected = true;
  }

  const longUrl = `https://example.com/${randomString(32)}`;

  const shortenRes = client.invoke(
    "shortener.Shortener/Shorten",
    { url: longUrl },
    { tags: { method: "Shorten" } }
  );
  check(shortenRes, {
    "Shorten OK": (r) => r && r.status === grpc.StatusOK && !!r.message.code,
  });
  const code = shortenRes.message.code;

  const resolveRes = client.invoke(
    "shortener.Shortener/Resolve",
    { code: code },
    { tags: { method: "Resolve" } }
  );
  check(resolveRes, {
    "Resolve OK": (r) => r && r.status === grpc.StatusOK && !!r.message.url,
  });
}

export function teardown() {
  if (connected) {
    client.close();
  }
}