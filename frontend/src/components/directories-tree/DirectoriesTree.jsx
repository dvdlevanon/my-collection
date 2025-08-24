import CheckIcon from '@mui/icons-material/Check';
import MenuIcon from '@mui/icons-material/Menu';
import { Chip, IconButton, Stack, Typography, useTheme } from '@mui/material';
import { useEffect, useRef, useState } from 'react';
import { Tree } from 'react-arborist';
import useResizeObserver from 'use-resize-observer';
import Client from '../../utils/client';
import PathUtil from '../../utils/path-util';
import CategoriesChooser from '../categories-chooser/CategoriesChooser';
import DirectorySettingsMenu from './DirectorySettingsMenu';
import { useTreeStore } from './TreeStore';

function DirectoriesTree() {
	const unknownCategoryId = -1;

	const { treeData, loadNode, updateNodeInfo, refreshNodeComplete } = useTreeStore();
	const treeRef = useRef();
	const { ref, width, height } = useResizeObserver();
	const theme = useTheme();
	const [hoverId, setHoverId] = useState('');
	const containerRef = useRef();
	const [dirMenuData, setDirMenuData] = useState(null);

	useEffect(() => {
		openNode(treeRef.current.root.children[0]);
	}, [treeRef.current]);

	const openNode = async (node) => {
		if (!node.data.isLoaded) {
			await loadNode(node.data.id, 1);
			node.open();
		} else if (node.data.children) {
			node.toggle();
		}
	};

	const setCategories = async (dir, categoryIds) => {
		let categories = [];
		for (let i = 0; i < categoryIds.length; i++) {
			categories.push({ id: categoryIds[i] });
		}

		await Client.setDirectoryCategories({ ...dir.dirinfo, tags: categories });
		updateNodeInfo(dir.id, {
			...dir.dirinfo,
			tags: categories,
		});
	};

	const renderDirInfo = (dir) => {
		if (!dir.dirinfo) {
			return;
		}

		if (dir.dirinfo.excluded) {
			return;
		}

		let selectedCategories = [unknownCategoryId];
		if (dir.dirinfo.tags) {
			selectedCategories = dir.dirinfo.tags.map((dir) => dir.id);
		}

		return (
			<Stack flexDirection="row" gap={theme.spacing(2)}>
				<Chip label="Included" icon={<CheckIcon sx={{ fontSize: theme.iconSize(1) }} />}></Chip>
				<CategoriesChooser
					multiselect={true}
					allowToCreate={true}
					placeholder="Not Specified"
					selectedIds={selectedCategories}
					setCategories={(categoriyIds) => {
						setCategories(dir, categoriyIds);
					}}
				/>
			</Stack>
		);
	};

	const renderRow = ({ node, style }) => {
		const icon = node.isOpen ? '📂' : '📁';

		return (
			<Stack
				flexDirection="row"
				style={style}
				sx={{
					cursor: 'pointer',
				}}
				alignItems="center"
				height="100%"
				gap={theme.spacing(2)}
				onMouseLeave={() => setHoverId('')}
				onMouseEnter={() => setHoverId(node.data.id)}
			>
				<Typography
					variant="h5"
					onClick={() => {
						node = treeRef.current.get(node.id);
						openNode(node);
					}}
				>
					{icon} {PathUtil.dirname(node.data.name)}
				</Typography>
				{renderDirInfo(node.data)}
				{hoverId == node.data.id && (
					<IconButton
						onClick={(e) => {
							setDirMenuData({
								dir: node.data,
								position: {
									left: e.clientX,
									top: e.clientY,
								},
							});
						}}
					>
						<MenuIcon />
					</IconButton>
				)}
			</Stack>
		);
	};

	return (
		<Stack
			width={'100%'}
			height={'100%'}
			flexGrow={1}
			overflow="hidden"
			backgroundColor="dark.lighter"
			className="parent"
			ref={(el) => {
				ref(el);
				containerRef.current = el;
			}}
			padding={theme.spacing(3)}
		>
			<Tree
				height={height}
				width={width}
				ref={treeRef}
				data={treeData}
				rowHeight={60}
				openByDefault={false}
				disableDrag
				disableEdit
				disableDrop
			>
				{renderRow}
			</Tree>
			{dirMenuData && (
				<DirectorySettingsMenu
					anchorEl={containerRef.current}
					anchorPosition={dirMenuData.position}
					dir={dirMenuData.dir}
					refreshDir={(dir) => refreshNodeComplete(dir.id)}
					onClose={() => {
						let nodeId = dirMenuData.dir.id;
						setDirMenuData(null);
						setHoverId('');
						let node = treeRef.current.get(nodeId);
						node.select();
					}}
				/>
			)}
		</Stack>
	);
}

export default DirectoriesTree;
