import { Tooltip } from '@mui/material';
import { Link } from 'react-router-dom';

function Item({ item }) {
	const getCover = () => {
		if (item.previews && item.previews.length > 0) {
			return "http://localhost:8080/storage/" + encodeURIComponent(item.previews[0].url)
		} else {
			return "empty"
		}
	}

    return (
		<Link to={"item/" + item.id} className="item">
			<img src={getCover()} alt='' />
			<Tooltip title={item.title} arrow followCursor >
			<span className="item-title">{item.title}</span>
			</Tooltip>
		</Link>
    )
}

export default Item
