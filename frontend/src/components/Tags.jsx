import { TextField } from '@mui/material';
import { useState } from 'react';
import Tag from './Tag';
import styles from './Tags.module.css';

function Tags({ tags, onTagSelected }) {
	let [searchTerm, setSearchTerm] = useState('');

	const onSearchTermChanged = (e) => {
		setSearchTerm(e.target.value);
	};

	const filterTags = (tags, searchTerm) => {
		if (!searchTerm) {
			return tags;
		}

		return tags.filter((tag) => {
			return tag.title.toLowerCase().includes(searchTerm.toLowerCase());
		});
	};

	return (
		<div className={styles.tags_holder}>
			<div className={styles.filters_holder}>
				<TextField
					autoFocus
					fullWidth
					label="Search for tags"
					type="search"
					onChange={(e) => onSearchTermChanged(e)}
				></TextField>
			</div>
			<div className={styles.tags}>
				{filterTags(tags, searchTerm).map((tag) => {
					return (
						<div key={tag.id}>
							<Tag tag={tag} onTagSelected={onTagSelected} />
						</div>
					);
				})}
			</div>
		</div>
	);
}

export default Tags;
