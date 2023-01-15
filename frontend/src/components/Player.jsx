import styles from './Player.module.css';

function Player({ item }) {
	return (
		<div className={styles.player}>
			<video muted controls width="100%" height="700px">
				<source src={'http://localhost:8080/stream/' + item.url} />
			</video>
		</div>
	);
}

export default Player;
