const openNewFactos = document.querySelector("[newfactos-open-modal]")
const closeNewFactos = document.querySelector("[newfactos-close-modal]")
const newFactosModal = document.querySelector("[newfactos-modal]")

openNewFactos.addEventListener("click", () => {
    console.log("Opening modal")
    newFactosModal.showModal()
})

closeNewFactos.addEventListener("click", () => {
    console.log("Closing Modal")
    newFactosModal.close()
})

