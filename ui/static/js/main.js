const openNewFactos = document.querySelectorAll("[newfactos-open-modal]");
const newFactosModal = document.querySelectorAll("[newfactos-modal]");
//const closeNewFactos = document.querySelectorAll("[newfactos-close-modal]")

const closeNewFactos = document.querySelectorAll(".newfactos-close-modal");
const exit = document.querySelectorAll(".cancel-close-modal");
const x_close = document.querySelectorAll(".x-close-modal");
console.log("Running JS script");

for (let i = 0; i < openNewFactos.length; i++) {
    openNewFactos.item(i).addEventListener("click", () => {
        console.log("Opening Modal");
        newFactosModal.item(i).showModal();
        
        closeNewFactos.item(i).addEventListener("click", () => {
            console.log("Closing Modal");
            newFactosModal.item(i).close();
        });

        exit.item(i).addEventListener("click", () => {
            console.log("Closing Modal");
            newFactosModal.item(i).close();
        });

        x_close.item(i).addEventListener("click", () => {
            console.log("Closing Modal");
            newFactosModal.item(i).close();
        });
    })

}
