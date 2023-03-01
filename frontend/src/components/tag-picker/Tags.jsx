import { Button, Link, Stack } from '@mui/material';
import { Box } from '@mui/system';
import { useState } from 'react';
import { useQuery } from 'react-query';
import { Link as RouterLink } from 'react-router-dom';
import seedrandom from 'seedrandom';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import AddTagDialog from '../dialogs/AddTagDialog';
import Tag from './Tag';
import TagsFilter from './TagsFilter';

function Tags({ tags, parentId, size, onTagSelected }) {
	const [addTagDialogOpened, setAddTagDialogOpened] = useState(false);
	const [searchTerm, setSearchTerm] = useState('');
	const [sortBy, setSortBy] = useState('random');
	const [selectedAnnotations, setSelectedAnnotations] = useState([]);
	const availableAnnotationsQuery = useQuery({
		queryKey: ReactQueryUtil.availableAnnotationsKey(parentId),
		queryFn: () => Client.getAvailableAnnotations(parentId),
		onSuccess: (availableAnnotations) => {
			setSelectedAnnotations(
				selectedAnnotations.filter((selected) => {
					return availableAnnotations.some((annotation) => selected.id == annotation.id);
				})
			);
		},
	});

	const filterTagsBySearch = (tags) => {
		let filteredTags = tags;

		if (searchTerm) {
			filteredTags = tags.filter((tag) => {
				return tag.title.toLowerCase().includes(searchTerm.toLowerCase());
			});
		}

		return filteredTags;
	};

	const filterTagsByAnnotations = (tags) => {
		return tags.filter((cur) => {
			if (selectedAnnotations.length == 0) {
				return true;
			}

			if (!cur.tags_annotations) {
				return false;
			}

			return cur.tags_annotations.some((tagAnnotation) => {
				return selectedAnnotations.some((annotation) => annotation.id == tagAnnotation.id);
			});
		});
	};

	const filterTags = () => {
		let filteredTags = filterTagsByAnnotations(filterTagsBySearch(tags));

		if (sortBy == 'random') {
			let epochDay = Math.floor(Date.now() / 1000 / 60 / 60 / 24);
			let randomTags = [];
			let rand = seedrandom(epochDay);

			for (let i = 0; i < filteredTags.length; i++) {
				let randomIndex = Math.floor(rand() * filteredTags.length);
				while (randomTags[randomIndex]) {
					randomIndex = Math.floor(rand() * filteredTags.length);
				}

				randomTags[randomIndex] = filteredTags[i];
			}

			return randomTags;
		} else if (sortBy == 'title-asc') {
			return filteredTags.sort((a, b) => (a.title > b.title ? 1 : a.title < b.title ? -1 : 0));
		} else if (sortBy == 'title-desc') {
			return filteredTags.sort((a, b) => (a.title > b.title ? -1 : a.title < b.title ? 1 : 0));
		} else {
			return filterTags;
		}
	};

	return (
		<Stack width={'100%'} height={'100%'} flexGrow={1} backgroundColor="dark.lighter">
			<TagsFilter
				parentId={parentId}
				setSearchTerm={setSearchTerm}
				annotations={(availableAnnotationsQuery.isSuccess && availableAnnotationsQuery.data) || []}
				setAddTagDialogOpened={setAddTagDialogOpened}
				selectedAnnotations={selectedAnnotations}
				setSelectedAnnotations={setSelectedAnnotations}
			/>
			<Box
				sx={{
					display: 'flex',
					flexDirection: 'row',
					padding: '10px',
					gap: '10px',
					flexWrap: 'wrap',
				}}
			>
				{filterTags().map((tag) => {
					return <Tag key={tag.id} tag={tag} size={size} onTagSelected={onTagSelected} />;
				})}
				{!TagsUtil.isDirectoriesCategory(parentId) && (
					<Tag key="add-tag" tag={{ id: -1 }} size={size} onTagSelected={() => setAddTagDialogOpened(true)} />
				)}
				{TagsUtil.isDirectoriesCategory(parentId) && tags.length == 0 && (
					<Stack
						direction="row"
						gap="10px"
						justifyContent="center"
						alignItems="center"
						sx={{
							width: '100%',
							height: '100%',
						}}
					>
						No diretories found
						<Link component={RouterLink} to="spa/manage-directories">
							<Button variant="outlined">Manage Directories</Button>
						</Link>
					</Stack>
				)}
			</Box>
			{addTagDialogOpened && (
				<AddTagDialog parentId={parentId} verb="Tag" onClose={() => setAddTagDialogOpened(false)} />
			)}
		</Stack>
	);
}

export default Tags;
