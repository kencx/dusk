document.body.addEventListener("openModal", function(event) {
	const modalDiv = document.querySelector("#modal-content")
	const modalDialog = modalDiv.querySelector("dialog")
	// const closeModalBtn = modalDiv.querySelector('button[aria-label="close"]#modal-close-btn')
	const cancelModalBtn = modalDiv.querySelector("#modal-cancel-btn")

	const closeModal = function() {
		modalDialog.removeAttribute("open")
	}

	// closeModalBtn.addEventListener("click", closeModal);
	cancelModalBtn.addEventListener("click", closeModal);

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
