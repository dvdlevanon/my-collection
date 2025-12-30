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
import { FS_NODE_DIR, FS_NODE_FILE, useTreeStore } from './TreeStore';

function DirectoriesTree() {
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

	const getParentsCategories = (dir) => {
		let result = [];
		while (dir) {
			if (dir.dirinfo && dir.dirinfo.tags) {
				result.push(...dir.dirinfo.tags.map((tag) => tag.id));
			}
			dir = dir.parent;
		}
		return result;
	};

	const renderDirInfo = (dir) => {
		if (!dir.dirinfo) {
			return;
		}

		if (dir.dirinfo.excluded) {
			return;
		}

		let parentCategories = getParentsCategories(dir.parent);
		let selectedCategories = [];
		if (dir.dirinfo.tags) {
			selectedCategories = dir.dirinfo.tags.map((tag) => tag.id);
		}

		return (
			<Stack flexDirection="row" gap={theme.spacing(2)}>
				<Chip label="Included" icon={<CheckIcon sx={{ fontSize: theme.iconSize(1) }} />}></Chip>
				<CategoriesChooser
					multiselect={true}
					allowToCreate={true}
					placeholder="Not Specified"
					immutableIds={parentCategories}
					immutablePlaceholder={' (inherited)'}
					selectedIds={selectedCategories}
					setCategories={(categoriyIds) => {
						setCategories(dir, categoriyIds);
					}}
				/>
			</Stack>
		);
	};

	const renderDirRow = (node, style) => {
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

	const getFileIcon = (filename) => {
		const ext = filename.split('.').pop()?.toLowerCase() || '';

		// Movies/Video
		if (['mp4', 'mkv', 'avi', 'mov', 'wmv', 'flv', 'webm', 'm4v'].includes(ext)) {
			return '🎬';
		}

		// Audio
		if (['mp3', 'wav', 'flac', 'aac', 'ogg', 'wma', 'm4a', 'opus'].includes(ext)) {
			return '🎵';
		}

		// Text/Documents
		if (['txt', 'md', 'doc', 'docx', 'pdf', 'rtf'].includes(ext)) {
			return '📄';
		}

		// Scripts/Code
		if (['js', 'py', 'go', 'sh', 'bash', 'ts', 'jsx', 'tsx', 'java', 'c', 'cpp', 'rs', 'rb'].includes(ext)) {
			return '📜';
		}

		// Images
		if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'svg', 'webp'].includes(ext)) {
			return '🖼️';
		}

		// Archives
		if (['zip', 'rar', '7z', 'tar', 'gz', 'bz2'].includes(ext)) {
			return '📦';
		}

		// Other
		return '📎';
	};

	const renderFileRow = (node, style) => {
		const icon = getFileIcon(node.data.name);

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
			>
				<Typography variant="h5">
					{icon} {PathUtil.dirname(node.data.name)}
				</Typography>
			</Stack>
		);
	};

	const renderRow = ({ node, style }) => {
		if (node.data.type == FS_NODE_DIR) {
			return renderDirRow(node, style);
		} else if (node.data.type == FS_NODE_FILE) {
			return renderFileRow(node, style);
		}
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
