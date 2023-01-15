import { Tooltip } from '@mui/material';
import { Link } from 'react-router-dom';
import styles from './Item.module.css';

function Item({ item }) {
	const getCover = () => {
		if (item.previews && item.previews.length > 0) {
			return 'http://localhost:8080/storage/' + encodeURIComponent(item.previews[0].url);
		} else {
			return 'empty';
		}
	};

	return (
		<Link className={styles.item} to={'item/' + item.id}>
			<img className={styles.image} src={getCover()} alt="" />
			<Tooltip title={item.title} arrow followCursor>
				<span className={styles.item_title}>{item.title}</span>
			</Tooltip>
		</Link>
	);
}

export default Item;
