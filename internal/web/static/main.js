const CONFIG = {
  debounceMs: Number("{{.DebounceMs}}") || 500,
  refreshMs: Number("{{.RefreshMs}}") || 5000,
};

let volumeTimer = null;
let pendingVolume = false;
let errorTimer = null;

function showError(msg) {
  const el = document.getElementById("error");
  el.textContent = msg;
  clearTimeout(errorTimer);
  errorTimer = setTimeout(() => {
    el.textContent = "";
  }, 3000);
}

function api(field, opts) {
  return fetch(`/api/${field}`, opts)
    .then((res) =>
      res.ok
        ? res
        : res.text().then((t) => {
            showError(t);
            return null;
          }),
    )
    .catch((e) => {
      showError(e.message);
      return null;
    });
}

function post(field, body) {
  return api(field, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  }).then(loadState);
}

function updateVolumeUI(v) {
  document.body.style.setProperty("--volume", `${v}%`);
  document
    .querySelectorAll(".volume-label")
    .forEach((el) => (el.textContent = v));
}

function load(field, cb) {
  return api(field).then((res) => res && res.json().then(cb));
}

function loadState() {
  return Promise.all([
    pendingVolume || load("volume", (data) => updateVolumeUI(data.volume)),
    load("mute", (data) => {
      document.querySelectorAll(".mute button").forEach((el) => {
        el.textContent = data.mute ? "🔇" : "🔊";
      });
    }),
  ]);
}

function setVolume(v) {
  v = Math.max(0, Math.min(100, Math.round(v)));
  pendingVolume = true;
  updateVolumeUI(v);
  clearTimeout(volumeTimer);
  volumeTimer = setTimeout(() => {
    post("volume", { volume: v }).finally(() => {
      pendingVolume = false;
    });
  }, CONFIG.debounceMs);
}

function volumeDrag(e) {
  if (e.buttons === 0 && e.type !== "pointerdown") return;
  const el = e.currentTarget;
  const rect = el.getBoundingClientRect();

  const margin = 0.05; // 5% margin at edges
  const mapWithMargin = (raw) => {
    const clamped = Math.max(margin, Math.min(1 - margin, raw));
    return (clamped - margin) / (1 - 2 * margin);
  };

  let v;
  if (el.classList.contains("vertical")) {
    const raw = (e.clientY - rect.top) / rect.height;
    v = 100 - mapWithMargin(raw) * 100;
  } else {
    const raw = (e.clientX - rect.left) / rect.width;
    v = mapWithMargin(raw) * 100;
  }
  setVolume(v);
}

for (const action of ["playpause", "next", "previous", "mute"]) {
  document.querySelectorAll(`.${action} button`).forEach((el) => {
    el.addEventListener("click", () => post(action));
  });
}

document.querySelectorAll(".volume-control").forEach((el) => {
  el.addEventListener("pointerdown", volumeDrag);
  el.addEventListener("pointermove", volumeDrag);
});

loadState();
setInterval(loadState, CONFIG.refreshMs);
