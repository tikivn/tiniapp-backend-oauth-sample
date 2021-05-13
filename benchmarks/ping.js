import http from "k6/http";
import { check, fail } from "k6";

import { Counter, Trend } from "k6/metrics";

let durationTrend = new Trend("_custom_duration_trend");
let successCounter = new Counter("_custom_success_count");

// export let options = {
//   vus: 100,
//   duration: "1m",
// };

export default function () {
  const response = http.get("http://localhost:8080/api/ping");

  if (response.status === 200) {
    durationTrend.add(response.timings.duration);
    successCounter.add(1);
  }

  if (
    !check(response, {
      "status code MUST be 200": (res) => res.status == 200,
    })
  ) {
    fail(
      `status code was *not* 200 but was ${response.status} with error ${response.error} and body ${response.body}`
    );
  }
}
