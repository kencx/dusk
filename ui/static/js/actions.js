document.body.addEventListener("openModal", function(event) {
	const modalDiv = document.querySelector("#modal-content")
	const modalDialog = modalDiv.querySelector("dialog")
	const cancelModalBtn = modalDiv.querySelector("#modal-cancel-btn")

	const closeModal = function() {
		modalDialog.removeAttribute("open")
	}

	cancelModalBtn.addEventListener("click", closeModal);
	modalDialog.addEventListener("click", closeModal);

	// close on Esc
	document.addEventListener("keydown", function(e) {
		if (e.key === "Escape" && modalDialog.hasAttribute("open")) {
			closeModal();
		}
	});
});

document.body.addEventListener("onToast", function(event) {
	const toasts = document.querySelectorAll(".toast");
	toasts.forEach(t => {
		const closeBtn = t.querySelector("button[aria-label='close']")
		t.addEventListener("click", () => t.remove());
		closeBtn.addEventListener("click", () => t.remove());
	});
});

document.body.addEventListener('htmx:afterRequest', function(evt) {
	const errorTarget = document.getElementById("htmx-alert")
	if (evt.detail.successful) {
		// Successful request, clear out alert
		errorTarget.setAttribute("hidden")
		errorTarget.innerText = "";
	} else if (evt.detail.failed && evt.detail.xhr) {
		// Server error with response contents, equivalent to htmx:responseError
		console.warn("Server error", evt.detail)
		const xhr = evt.detail.xhr;
		errorTarget.innerText = `Unexpected server error: ${xhr.status} - ${xhr.statusText}`;
		errorTarget.removeAttribute("hidden");
	} else {
		// Unspecified failure, usually caused by network error
		console.error("Unexpected htmx error", evt.detail)
		errorTarget.innerText = "Unexpected error, check your connection and try to refresh the page.";
		errorTarget.removeAttribute("hidden");
	}
});
