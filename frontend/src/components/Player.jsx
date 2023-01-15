import Client from '../network/client';
import styles from './Player.module.css';

function Player({ item }) {
	return (
		<div className={styles.player}>
			<video muted controls width="100%" height="700px">
				<source src={Client.buildStreamUrl(item.url)} />
			</video>
		</div>
	);
}

export default Player;
