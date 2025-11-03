import React from "react";
import styles from "./Home.module.css";
import HeroImg from "../../assets/images/futbol.jpg";


export const Home: React.FC = () => {
	return(
		<div className={styles.main}>
      <div className={styles.text_container}>
        <h1>FACT<span style={{"color": "var(--secondary)"}}>O</span>S</h1>
        <p>Predict results for upcoming games.<br/>Show your friends who knows football best.</p>
        
      </div>
		</div>
	)
}