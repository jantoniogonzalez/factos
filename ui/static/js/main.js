const openNewFactos = document.querySelector("[newfactos-open-modal]")
const newFactosModal = document.querySelector("[newfactos-modal]")
//const closeNewFactos = document.querySelectorAll("[newfactos-close-modal]")

openNewFactos.addEventListener("click", () => {
    console.log("Opening modal")
    newFactosModal.showModal()
})

const closeNewFactos = document.querySelectorAll(".newfactos-close-modal")

closeNewFactos.addEventListener("click", () => {
    console.log("Closing Modal")
    newFactosModal.close()
})


