document.body.addEventListener("click", function (e) {
  if (e.target.classList.contains("expander")) {
    const ul = e.target.nextElementSibling.nextElementSibling;
    if (ul.classList.contains("collapsed")) {
      ul.classList.remove("collapsed");
      e.target.textContent = e.target.textContent.replace("▶", "▼");
    } else {
      ul.classList.add("collapsed");
      e.target.textContent = e.target.textContent.replace("▼", "▶");
    }
  }
});
