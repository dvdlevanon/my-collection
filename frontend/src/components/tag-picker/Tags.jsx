import { Button, Divider, Link, Stack } from '@mui/material';
import { Box } from '@mui/system';
import { useEffect, useState } from 'react';
import { useQuery } from 'react-query';
import { Link as RouterLink } from 'react-router-dom';
import seedrandom from 'seedrandom';
import AspectRatioUtil from '../../utils/aspect-ratio-util';
import Client from '../../utils/client';
import ReactQueryUtil from '../../utils/react-query-util';
import TagsUtil from '../../utils/tags-util';
import AddTagDialog from '../dialogs/AddTagDialog';
import Tag from './Tag';
import TagsTopBar from './TagsTopBar';

function Tags({ origin, tags, tits, parent, initialTagSize, tagLinkBuilder, onTagClicked }) {
	const [addTagDialogOpened, setAddTagDialogOpened] = useState(false);
	const [searchTerm, setSearchTerm] = useState('');
	const [sortBy, setSortBy] = useState(parent.default_sorting);
	const [tit, setTit] = useState(tits[0]);
	const [prefixFilter, setPrefixFilter] = useState('');
	const [selectedAnnotations, setSelectedAnnotations] = useState([]);
	const [tagSize, setTagSize] = useState(initialTagSize);
	const availableAnnotationsQuery = useQuery({
		queryKey: ReactQueryUtil.availableAnnotationsKey(parent.id),
		queryFn: () => Client.getAvailableAnnotations(parent.id),
		onSuccess: (availableAnnotations) => {
			setSelectedAnnotations(
				selectedAnnotations.filter((selected) => {
					return (
						selected.id == 'none' ||
						selected.id == 'no-image' ||
						availableAnnotations.some((annotation) => selected.id == annotation.id)
					);
				})
			);

			if (TagsUtil.isDailymixCategory(parent.id)) {
				var today = new Date();
				var month = today.toLocaleString('default', { month: 'short' });
				var year = today.getFullYear();
				var defaultTagAnnotation = month + '-' + year;
				setSelectedAnnotations(availableAnnotations.filter((cur) => cur.title == defaultTagAnnotation));
			}

			let lastSelectedAnnotations = localStorage.getItem(buildStorageKey('selected-annotations'));
			if (lastSelectedAnnotations) {
				lastSelectedAnnotations = lastSelectedAnnotations.split(',').map((cur) => cur.trim());
				setSelectedAnnotations(
					availableAnnotations.filter((availableAnnoation) => {
						return lastSelectedAnnotations.some((cur) => cur == availableAnnoation.id);
					})
				);

				if (lastSelectedAnnotations.some((cur) => cur == 'no-image'))
					setSelectedAnnotations((selectedAnnotations) => {
						let result = selectedAnnotations;
						result.push({
							id: 'no-image',
							title: 'No Image',
						});

						return result;
					});
			}
		},
	});

	useEffect(() => {
		let lastTit = localStorage.getItem(buildStorageKey('tit'));
		if (lastTit) {
			setTit(tits.find((cur) => cur.id == lastTit));
		}

		let lastTagSize = localStorage.getItem(buildStorageKey('tag-size'));
		if (lastTagSize) {
			setTagSize(parseInt(lastTagSize));
		}

		let lastSortBy = localStorage.getItem(buildStorageKey('sort-by'));
		if (lastSortBy) {
			setSortBy(lastSortBy);
		} else {
			setSortBy(parent.default_sorting);
		}

		let lastPrefixFilter = localStorage.getItem(buildStorageKey('prefix-filter'));
		if (lastPrefixFilter) {
			setPrefixFilter(lastPrefixFilter);
		}
	}, [origin, parent]);

	const buildStorageKey = (name) => {
		return 'tags-' + origin + '-' + parent.id + '-' + name;
	};

	const filterTagsBySearch = (tags) => {
		let filteredTags = tags;

		if (searchTerm) {
			filteredTags = tags.filter((tag) => {
				return tag.title.toLowerCase().includes(searchTerm.toLowerCase());
			});
		}

		return filteredTags;
	};

	const filterByPrefix = (tags) => {
		let filteredTags = tags;

		if (prefixFilter) {
			filteredTags = tags.filter((tag) => {
				return tag.title.toLowerCase().startsWith(prefixFilter.toLowerCase());
			});
		}

		return filteredTags;
	};

	const filterTagsByAnnotations = (tags) => {
		return tags.filter((cur) => {
			if (selectedAnnotations.length == 0) {
				return true;
			}

			let isNoImageSelected = selectedAnnotations.some((annotation) => annotation.id == 'no-image');
			if (isNoImageSelected) {
				return !TagsUtil.hasImage(cur);
			}

			if (!cur.tags_annotations) {
				let isNoneSelected = selectedAnnotations.some((annotation) => annotation.id == 'none');
				return isNoneSelected;
			}

			return cur.tags_annotations.some((tagAnnotation) => {
				return selectedAnnotations.some((annotation) => annotation.id == tagAnnotation.id);
			});
		});
	};

	const filterTags = () => {
		let filteredTags = filterByPrefix(filterTagsByAnnotations(filterTagsBySearch(tags)));

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
		} else if (sortBy == 'items-count') {
			return filteredTags.sort((a, b) =>
				TagsUtil.itemsCount(a) > TagsUtil.itemsCount(b)
					? -1
					: TagsUtil.itemsCount(a) < TagsUtil.itemsCount(b)
					? 1
					: 0
			);
		} else {
			return filteredTags;
		}
	};

	const getAvailableAnnotations = (availableAnnotations) => {
		if (availableAnnotations.length == 0) {
			return availableAnnotations;
		}

		if (TagsUtil.isSpecialCategory(parent.id)) {
			return availableAnnotations;
		}

		return [
			{
				id: 'none',
				title: 'None',
			},
			{
				id: 'no-image',
				title: 'No Image',
			},
			...availableAnnotations,
		];
	};

	const calculateTagSize = () => {
		let result = {};
		if (parent.display_style === 'portrait') {
			result = { width: AspectRatioUtil.calcHeight(tagSize, AspectRatioUtil.asepctRatio16_9), height: tagSize };
		} else if (parent.display_style === 'landscape') {
			result = { width: tagSize, height: AspectRatioUtil.calcHeight(tagSize, AspectRatioUtil.asepctRatio16_9) };
		} else if (parent.display_style == 'icon') {
			result = { width: tagSize / 4, height: tagSize / 4 };
		} else if (parent.display_style === 'banner') {
			result = { width: tagSize, height: AspectRatioUtil.calcHeight(tagSize, AspectRatioUtil.asepctRatio16_9) };
		} else if (parent.display_style === 'chip') {
			result = { width: 'auto', height: 'auto' };
		} else {
			console.log('unsupported display style ' + parent.display_style);
			result = { width: 'auto', height: 'auto' };
		}

		return result;
	};

	return (
		<Stack width={'100%'} height={'100%'} flexGrow={1} backgroundColor="dark.lighter">
			<TagsTopBar
				parentId={parent.id}
				setSearchTerm={setSearchTerm}
				annotations={
					(availableAnnotationsQuery.isSuccess && getAvailableAnnotations(availableAnnotationsQuery.data)) ||
					[]
				}
				setAddTagDialogOpened={setAddTagDialogOpened}
				selectedAnnotations={selectedAnnotations}
				setSelectedAnnotations={(selectedAnnotations) => {
					setSelectedAnnotations(selectedAnnotations);
					localStorage.setItem(
						buildStorageKey('selected-annotations'),
						selectedAnnotations.map((cur) => cur.id)
					);
				}}
				tits={tits}
				tit={tit}
				setTit={(tit) => {
					localStorage.setItem(buildStorageKey('tit'), tit.id);
					setTit(tit);
				}}
				sortBy={sortBy}
				setSortBy={(sortBy) => {
					localStorage.setItem(buildStorageKey('sort-by'), sortBy);
					setSortBy(sortBy);
				}}
				prefixFilter={prefixFilter}
				setPrefixFilter={(prefixFilter) => {
					localStorage.setItem(buildStorageKey('prefix-filter'), prefixFilter);
					setPrefixFilter(prefixFilter);
				}}
				tagSize={tagSize}
				onZoomChanged={(offset) => {
					let newTagSize = parseInt(tagSize) + parseInt(offset);
					setTagSize(newTagSize);
					localStorage.setItem(buildStorageKey('tag-size'), newTagSize);
				}}
			/>
			<Divider />
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
					return (
						<Tag
							key={tag.id}
							tag={tag}
							parent={parent}
							tagDimension={calculateTagSize()}
							selectedTit={tit}
							tagLinkBuilder={tagLinkBuilder}
							onTagClicked={onTagClicked}
						/>
					);
				})}
				{TagsUtil.isSpecialCategory(parent.id) && tags.length == 0 && (
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
				<AddTagDialog parentId={parent.id} verb="Tag" onClose={() => setAddTagDialogOpened(false)} />
			)}
		</Stack>
	);
}

export default Tags;
