import './App.css'
import gsap from "gsap";
import { useGSAP } from '@gsap/react';
import { Home } from './pages/Home/Home'

gsap.registerPlugin(useGSAP);

function App() {
  return (
    <>
      <Home />
    </>
  )
}

export default App
