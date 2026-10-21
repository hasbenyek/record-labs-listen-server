# Record Labs — Listen Together server for Render (free web service).
# MetrolistGroup's metroserver (GPL-3.0), the exact build tested with Record Labs
# (commit 9da6f9d, 2026-09-28), plus a policy that lets Record Labs host rooms.
FROM ghcr.io/metrolistgroup/metroserver@sha256:9655dfbc258c2f27bad81dabe60af256da768fa7cfd7040fbff7231a4e6c1265

COPY --chown=1000:1000 ua_policy.json /app/ua_policy.json
ENV UA_POLICY_FILE=/app/ua_policy.json
# Render sets PORT itself; metroserver listens on it.
