let currentSession = null;
let appInfo = null;

const $ = (id) => document.getElementById(id);
const esc = (s) => String(s ?? "").replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/"/g, "&quot;");

async function jsonFetch(url, opts) {
  const res = await fetch(url, opts);
  if (res.status === 204) return null;
  const body = await res.json();
  if (!res.ok || body.error) {
    const err = new Error((body.error && body.error.message) || res.statusText);
    err.status = res.status;
    err.code = body.error && body.error.code;
    throw err;
  }
  return body.data;
}

async function patchRows(sessionId, tableName, assumedRowCount) {
  return jsonFetch(
    `/api/v1/analysis-sessions/${sessionId}/tables/${encodeURIComponent(tableName)}`,
    {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ assumed_row_count: Number(assumedRowCount) }),
    }
  );
}

async function patchGrowth(sessionId, tableName, rowsPerPeriod, period, horizon) {
  return jsonFetch(
    `/api/v1/analysis-sessions/${sessionId}/tables/${encodeURIComponent(tableName)}/growth`,
    {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        rows_per_period: Number(rowsPerPeriod),
        period,
        horizon: Number(horizon),
      }),
    }
  );
}

function exportURL(sessionId, format) {
  return `/api/v1/analysis-sessions/${sessionId}/export?format=${format}`;
}

const fmt = (n) => {
  if (n == null || n === "") return "—";
  const bytes = Number(n);
  if (!Number.isFinite(bytes)) return "—";
  const unit = 1024;
  if (Math.abs(bytes) < unit) return `${Math.trunc(bytes)} B`;
  const units = ["KB", "MB", "GB", "TB"];
  let div = unit;
  let exp = 0;
  for (let v = bytes / unit; v >= unit && exp < units.length - 1; v /= unit) {
    div *= unit;
    exp++;
  }
  return `${(bytes / div).toFixed(1)} ${units[exp]}`;
};

const totalLabel = (d) => d.estimated_total_human || fmt(d.estimated_total_bytes);
const projectedLabel = (d) =>
  d.projected_total_human || (d.projected_total_bytes != null ? fmt(d.projected_total_bytes) : null);

function colAttrs(col) {
  const bits = [];
  if (col.length != null) bits.push(`len ${col.length}`);
  if (col.precision != null) {
    bits.push(col.scale != null ? `p${col.precision},s${col.scale}` : `p${col.precision}`);
  }
  if (col.assumed_avg_length != null) bits.push(`avg ${col.assumed_avg_length}`);
  return bits.length ? bits.join(" · ") : "—";
}

function columnsHTML(cols) {
  if (!cols || !cols.length) {
    return `<div class="inner-card"><h4>Columns</h4><p class="meta">None</p></div>`;
  }
  let h = `<div class="inner-card"><h4>Columns</h4><table><thead><tr>
    <th>Name</th><th>Type</th><th>Normalized</th><th>Attrs</th><th>Est. size</th>
  </tr></thead><tbody>`;
  for (const col of cols) {
    h += `<tr>
      <td>${esc(col.name)}</td>
      <td>${esc(col.original_type || "?")}</td>
      <td>${esc(col.normalized_type || "—")}</td>
      <td>${esc(colAttrs(col))}</td>
      <td>${fmt(col.estimated_bytes)}</td>
    </tr>`;
  }
  return h + `</tbody></table></div>`;
}

function indexesHTML(indexes) {
  if (!indexes || !indexes.length) {
    return `<div class="inner-card"><h4>Indexes</h4><p class="meta">None detected</p></div>`;
  }
  let h = `<div class="inner-card"><h4>Indexes</h4><table><thead><tr>
    <th>Name</th><th>Kind</th><th>Columns</th><th>Est. size</th>
  </tr></thead><tbody>`;
  for (const ix of indexes) {
    const cols = Array.isArray(ix.columns) ? ix.columns.join(", ") : "";
    h += `<tr>
      <td>${esc(ix.name)}</td>
      <td>${esc(ix.kind)}</td>
      <td>${esc(cols)}</td>
      <td>${fmt(ix.estimated_bytes)}</td>
    </tr>`;
  }
  return h + `</tbody></table></div>`;
}

function growthHTML(t) {
  const rpp = t.growth_rows_per_period ?? 100;
  const period = t.growth_period || "day";
  const horizon = t.growth_horizon || 30;
  const attr = esc(t.name);
  const hasProj = t.projected_row_count != null || t.projected_table_bytes != null;
  const projRows = hasProj ? Number(t.projected_row_count).toLocaleString() : "—";
  const projBytes = hasProj ? fmt(t.projected_table_bytes) : "—";
  return `<div class="inner-card growth-card">
    <h4>Growth</h4>
    <div class="growth-layout">
      <div class="growth">
        <label>Rows / period
          <input class="field" data-growth-rpp="${attr}" type="number" min="1" value="${rpp}" />
        </label>
        <label>Period
          <select class="field" data-growth-period="${attr}">
            <option value="hour"${period === "hour" ? " selected" : ""}>hour</option>
            <option value="day"${period === "day" ? " selected" : ""}>day</option>
            <option value="month"${period === "month" ? " selected" : ""}>month</option>
          </select>
        </label>
        <label>Horizon
          <input class="field" data-growth-horizon="${attr}" type="number" min="1" value="${horizon}" />
        </label>
        <button type="button" class="btn-compact" data-apply-growth="${attr}">Apply</button>
      </div>
      <div class="projected-kpi${hasProj ? "" : " is-empty"}">
        <span class="label">Projected</span>
        <div class="projected-values">
          <span class="value">${projRows} <small>rows</small></span>
          <span class="value accent">${projBytes}</span>
        </div>
      </div>
    </div>
  </div>`;
}

function calcHTML(c) {
  return `<div class="inner-card">
    <h4>Calculation</h4>
    <div class="kpi-strip">
      <div class="kpi"><span class="label">Payload</span><span class="value">${fmt(c.column_payload_bytes)}</span></div>
      <div class="kpi"><span class="label">Header</span><span class="value">${fmt(c.row_header_bytes)}</span></div>
      <div class="kpi"><span class="label">Null bitmap</span><span class="value">${fmt(c.null_bitmap_bytes)}</span></div>
      <div class="kpi"><span class="label">Indexes</span><span class="value">${fmt(c.index_bytes)}</span></div>
    </div>
  </div>`;
}

async function loadInfo() {
  try {
    appInfo = await jsonFetch("/api/v1/info");
    const persist = appInfo.persistence ? "on" : "off";
    const ai = appInfo.ai_insights ? "on" : "off";
    $("status").textContent =
      `v${appInfo.version} · ${appInfo.bind} · local-only · persistence ${persist} · AI insights ${ai} · external_required ${appInfo.external_required}`;
  } catch {
    $("status").textContent = "Could not load /api/v1/info";
  }
}

async function refreshHistory() {
  const el = $("history");
  try {
    const items = await jsonFetch("/api/v1/analysis-sessions");
    if (!items || !items.length) {
      el.innerHTML = `<p class="meta" style="margin:0">No sessions yet. Run Analyze to create one.</p>`;
      return;
    }
    let html = `<div class="session-list">`;
    for (const s of items) {
      html += `<div class="session-card">
        <code class="id">${esc(s.id)}</code>
        <div class="meta-line">${esc(s.engine)} · ${s.table_count ?? 0} tables</div>
        <div class="actions">
          <button type="button" class="ghost" data-open-session="${esc(s.id)}">Open</button>
          <button type="button" class="ghost" data-delete-session="${esc(s.id)}">Delete</button>
        </div>
      </div>`;
    }
    el.innerHTML = html + `</div>`;
    el.querySelectorAll("[data-open-session]").forEach((btn) => {
      btn.onclick = () => openSession(btn.getAttribute("data-open-session"));
    });
    el.querySelectorAll("[data-delete-session]").forEach((btn) => {
      btn.onclick = () => deleteSession(btn.getAttribute("data-delete-session"));
    });
  } catch (e) {
    el.textContent = e.message || String(e);
  }
}

async function openSession(id) {
  try {
    $("err").textContent = "";
    const d = await jsonFetch(`/api/v1/analysis-sessions/${encodeURIComponent(id)}`);
    if (d.engine) $("engine").value = d.engine;
    if (d.ddl_text != null) $("ddl").value = d.ddl_text;
    await render(d);
  } catch (e) {
    $("err").textContent = e.message || String(e);
  }
}

async function deleteSession(id) {
  try {
    $("err").textContent = "";
    await jsonFetch(`/api/v1/analysis-sessions/${encodeURIComponent(id)}`, { method: "DELETE" });
    if (currentSession && currentSession.id === id) {
      currentSession = null;
      $("out").innerHTML = "";
    }
    await refreshHistory();
  } catch (e) {
    $("err").textContent = e.message || String(e);
  }
}

async function refreshSettings() {
  const el = $("settings");
  try {
    const data = await jsonFetch("/api/v1/settings");
    const list = (data && data.settings) || [];
    let html = "";
    if (!list.length) {
      html += `<p class="meta">No settings stored yet.</p>`;
    } else {
      html += `<table><thead><tr><th>Key</th><th>Value</th><th>Type</th></tr></thead><tbody>`;
      for (const s of list) {
        html += `<tr><td>${esc(s.key)}</td><td>${esc(s.value)}</td><td>${esc(s.value_type)}</td></tr>`;
      }
      html += `</tbody></table>`;
    }
    html += `<div style="display:grid;gap:0.55rem;margin-top:0.75rem">
      <label>Key <input class="field" id="setting-key" style="width:100%;display:block" /></label>
      <label>Value <input class="field" id="setting-value" style="width:100%;display:block" /></label>
      <label>Type
        <select class="field" id="setting-type" style="width:100%;display:block">
          <option value="string">string</option>
          <option value="number">number</option>
          <option value="bool">bool</option>
        </select>
      </label>
      <button type="button" id="setting-save">Save</button>
    </div>`;
    el.innerHTML = html;
    $("setting-save").onclick = async () => {
      try {
        $("err").textContent = "";
        await jsonFetch("/api/v1/settings", {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            key: $("setting-key").value,
            value: $("setting-value").value,
            value_type: $("setting-type").value,
          }),
        });
        await refreshSettings();
      } catch (e) {
        $("err").textContent = e.message || String(e);
      }
    };
  } catch (e) {
    if (e.status === 404) {
      el.innerHTML = `<p class="meta" style="margin:0">Persistence is off. Start with <code>tabyte serve --persist</code> to edit settings and keep history across restarts.</p>`;
      return;
    }
    el.textContent = e.message || String(e);
  }
}

async function loadInsights(sessionId) {
  const host = document.getElementById("insights");
  if (!host) return;
  try {
    const data = await jsonFetch(`/api/v1/analysis-sessions/${encodeURIComponent(sessionId)}/insights`);
    if (!data.enabled) {
      host.innerHTML = `<p class="meta" style="margin:0">AI insights extension is disabled (no external calls).</p>`;
      return;
    }
    const items = data.insights || [];
    if (!items.length) {
      host.innerHTML = `<p class="meta" style="margin:0">No insights returned.</p>`;
      return;
    }
    host.innerHTML =
      `<ul>` +
      items
        .map(
          (i) =>
            `<li><strong>${esc(i.category || i.provider || "insight")}</strong> — ${esc(i.text)}</li>`
        )
        .join("") +
      `</ul>`;
  } catch (e) {
    host.textContent = e.message || String(e);
  }
}

function tableCardHTML(t) {
  const c = t.calculation || {};
  const attr = esc(t.name);
  return `<article class="table-card" data-table="${attr}">
    <div class="table-card-head">
      <div>
        <h3>${attr}</h3>
        <p class="subline">${t.column_count ?? 0} columns · ${t.index_count ?? (t.indexes || []).length} indexes</p>
      </div>
      <div class="stat volume-stat">
        <span class="label">Volume</span>
        <span class="value">${fmt(t.estimated_table_bytes)}</span>
      </div>
    </div>
    <div class="metrics">
      <div class="metric"><span class="label">Row size</span><span class="value">${fmt(t.estimated_row_bytes)}</span></div>
      <div class="metric">
        <span class="label">Assumed rows</span>
        <div class="rows-control">
          <input class="field" data-rows="${attr}" type="number" min="0" value="${t.assumed_row_count}" />
          <button type="button" data-apply-rows="${attr}">Apply</button>
        </div>
      </div>
    </div>
    ${calcHTML(c)}
    <div class="card-grid two" style="margin-top:1rem">
      ${columnsHTML(t.columns)}
      ${indexesHTML(t.indexes)}
    </div>
    <div style="margin-top:1rem">${growthHTML(t)}</div>
  </article>`;
}

async function render(d) {
  currentSession = d;
  const tables = d.tables || [];
  const proj = projectedLabel(d);

  let html = `<section class="summary-card">
    <p class="disclaimer">Estimates are not physical measurements on a live database.</p>
    <div class="summary-stats">
      <div class="stat"><span class="label">Total</span><span class="value">${totalLabel(d)}</span></div>
      <div class="stat"><span class="label">Engine</span><span class="value sm">${esc(d.engine)}</span></div>
      <div class="stat"><span class="label">Tables</span><span class="value">${tables.length}</span></div>
      <div class="stat"><span class="label">Warnings</span><span class="value">${d.warning_count ?? 0}</span></div>
      <div class="stat"><span class="label">Signals</span><span class="value">${d.signal_count ?? 0}</span></div>
      ${proj ? `<div class="stat"><span class="label">Projected</span><span class="value">${proj}</span></div>` : ""}
      <div class="stat"><span class="label">Session</span><span class="value sm"><code>${esc(d.id)}</code></span></div>
    </div>
    <p class="export-links">
      <a href="${exportURL(d.id, "json")}" download>Export JSON</a> ·
      <a href="${exportURL(d.id, "csv")}" download>Export CSV</a>
    </p>
  </section>`;

  for (const t of tables) {
    html += tableCardHTML(t);
  }

  if ((d.warnings || []).length) {
    html += `<section class="alerts-card"><h3>Warnings</h3><ul>` +
      d.warnings
        .map(
          (w) =>
            `<li><strong>${esc(w.code)}</strong> ${esc(w.table)}${w.column ? "." + esc(w.column) : ""} — ${esc(w.message)}</li>`
        )
        .join("") +
      `</ul></section>`;
  }
  if ((d.signals || []).length) {
    html += `<section class="alerts-card"><h3>Signals</h3><ul>` +
      d.signals
        .map((s) => `<li><strong>${esc(s.code)}</strong> ${esc(s.table)} — ${esc(s.message)}</li>`)
        .join("") +
      `</ul></section>`;
  }

  html += `<section class="insights-card"><h3>Insights</h3><div id="insights" class="meta">Loading…</div></section>`;
  $("out").innerHTML = html;

  const runSafe = async (fn) => {
    try {
      $("err").textContent = "";
      await render(await fn());
      await refreshHistory();
    } catch (e) {
      $("err").textContent = e.message || String(e);
    }
  };

  $("out").querySelectorAll("[data-apply-rows]").forEach((btn) => {
    btn.onclick = () => {
      const name = btn.getAttribute("data-apply-rows");
      const card = btn.closest(".table-card");
      const input = card.querySelector(`[data-rows="${CSS.escape(name)}"]`);
      runSafe(() => patchRows(d.id, name, input.value));
    };
  });
  $("out").querySelectorAll("[data-apply-growth]").forEach((btn) => {
    btn.onclick = () => {
      const name = btn.getAttribute("data-apply-growth");
      const root = btn.closest(".inner-card");
      const rpp = root.querySelector("[data-growth-rpp]").value;
      const period = root.querySelector("[data-growth-period]").value;
      const horizon = root.querySelector("[data-growth-horizon]").value;
      runSafe(() => patchGrowth(d.id, name, rpp, period, horizon));
    };
  });

  await loadInsights(d.id);
}

$("refresh-history").onclick = () => refreshHistory();

$("run").onclick = async () => {
  $("err").textContent = "";
  $("out").innerHTML = "";
  $("run").disabled = true;
  try {
    const d = await jsonFetch("/api/v1/analysis-sessions", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        engine: $("engine").value,
        source_name: "paste.sql",
        ddl_text: $("ddl").value,
      }),
    });
    await render(d);
    await refreshHistory();
  } catch (e) {
    $("err").textContent = e.message || String(e);
  } finally {
    $("run").disabled = false;
  }
};

loadInfo();
refreshHistory();
refreshSettings();
