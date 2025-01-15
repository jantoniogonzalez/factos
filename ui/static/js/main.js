const openNewFactos = document.querySelector("[newfactos-open-modal]");
const newFactosModal = document.querySelector("[newfactos-modal]");
//const closeNewFactos = document.querySelectorAll("[newfactos-close-modal]")

openNewFactos.addEventListener("click", () => {
    console.log("Opening modal");
    newFactosModal.showModal();
});

const closeNewFactos = document.querySelector(".newfactos-close-modal");
const exit = document.querySelector(".cancel-close-modal");
const x_close = document.querySelector(".x-close-modal");

closeNewFactos.addEventListener("click", () => {
    console.log("Closing Modal");
    newFactosModal.close();
});

exit.addEventListener("click", () => {
    console.log("Closing Modal");
    newFactosModal.close();
});

x_close.addEventListener("click", () => {
    console.log("Closing Modal");
    newFactosModal.close();
});
