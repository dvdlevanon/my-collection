import Client from '../network/client';
import styles from './Player.module.css';

function Player({ item, isPreview }) {
	return (
		<div className={styles.player + ' ' + (isPreview && styles.preview)}>
			<video
				playsInline
				muted
				autoPlay={isPreview}
				loop={isPreview}
				controls={!isPreview}
				width="100%"
				height={isPreview ? '100%' : '700px'}
			>
				<source src={Client.buildFileUrl(isPreview ? item.previewUrl : item.url)} />
			</video>
		</div>
	);
}

export default Player;
