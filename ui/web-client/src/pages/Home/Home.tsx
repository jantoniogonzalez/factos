import React, { useRef } from "react";
import styles from "./Home.module.css";
import gsap from "gsap";
import { useGSAP } from '@gsap/react';


export const Home: React.FC = () => {
  const h1Ref = useRef<HTMLHeadingElement>(null);
  const divRef = useRef<HTMLDivElement>(null);
  const pRef = useRef<HTMLParagraphElement>(null);

  useGSAP(
    () => {
      const t1 = gsap.timeline();
      
      t1.from(h1Ref.current, {x: "-500", opacity: 0, duration: 1.5})
        .from(pRef.current, {y: "100", opacity: 0, duration: 0.5}, "-=0.75");
    },
    {scope: divRef}
  )

	return(
		<div className={styles.main}>
      <div ref={divRef} className={styles.text_container}>
        <h1 ref={h1Ref}>FACT<span>O</span>S</h1>
        <p ref={pRef}>Predict results for upcoming games.<br/>Show your friends who knows football best.</p>
      </div>
		</div>
	)
}