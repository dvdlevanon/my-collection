import { TextField } from '@mui/material';
import { useState } from 'react';
import Tag from './Tag';
import styles from './Tags.module.css';

function Tags({ tags, size, onTagSelected }) {
	let [searchTerm, setSearchTerm] = useState('');

	const onSearchTermChanged = (e) => {
		setSearchTerm(e.target.value);
	};

	const filterTags = (tags, searchTerm) => {
		let filteredTags = tags;

		if (searchTerm) {
			filteredTags = tags.filter((tag) => {
				return tag.title.toLowerCase().includes(searchTerm.toLowerCase());
			});
		}

		let shuffeledTags = filteredTags.sort(() => (Math.random() > 0.5 ? 1 : -1));
		return shuffeledTags;
	};

	return (
		<div className={styles.tags_holder}>
			<div className={styles.filters_holder}>
				<TextField
					variant="outlined"
					autoFocus
					fullWidth
					label="Search..."
					type="search"
					onChange={(e) => onSearchTermChanged(e)}
				></TextField>
			</div>
			<div className={styles.tags}>
				{filterTags(tags, searchTerm).map((tag) => {
					return (
						<div key={tag.id}>
							<Tag tag={tag} size={size} onTagSelected={onTagSelected} />
						</div>
					);
				})}
			</div>
		</div>
	);
}

export default Tags;
