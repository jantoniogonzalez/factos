const openNewFactos = document.querySelectorAll("[newfactos-open-modal]");
const newFactosModal = document.querySelectorAll("[newfactos-modal]");
//const closeNewFactos = document.querySelectorAll("[newfactos-close-modal]")
console.log("Running JS script");

for (let i = 0; i < openNewFactos.length; i++) {
    openNewFactos.item(i).addEventListener("click", () => {
        console.log("Opening Modal");
        newFactosModal.item(i).showModal();
    })
}
  
  

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
