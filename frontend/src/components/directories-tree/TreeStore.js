import { produce } from 'immer';
import { create } from 'zustand';
import Client from '../../utils/client';

const FS_NODE_DIR = 1;
const ROOT_DIR = '<root>';

const findNode = (nodes, nodeId) => {
	for (const node of nodes) {
		if (node.id === nodeId) return node;
		if (node.children) {
			const found = findNode(node.children, nodeId);
			if (found) return found;
		}
	}
	return null;
};

export const useTreeStore = create((set, get) => ({
	treeData: [
		{
			id: ROOT_DIR,
			name: 'root',
			type: FS_NODE_DIR,
			dirinfo: {},
			children: [],
			parent: null,
			isLoaded: false,
		},
	],

	loadNode: async (nodeId, depth = 1) => {
		const loadedDir = await Client.getFsDir(nodeId, depth);
		const { treeData } = get();

		set({
			treeData: produce(treeData, (draft) => {
				const node = findNode(draft, nodeId);
				if (!node) return;

				node.dirinfo = loadedDir.dirinfo;
				node.name = nodeId === ROOT_DIR ? loadedDir.path : node.name;
				node.isLoaded = true;

				if (loadedDir.children) {
					const dirChildren = loadedDir.children.filter((child) => child.type === FS_NODE_DIR);
					node.children = dirChildren.map((child) => ({
						id: child.path,
						name: child.path,
						dirinfo: child.dirinfo,
						children: [],
						parent: node,
					}));
				}
			}),
		});
	},

	updateNodeInfo: (nodeId, dirinfo) => {
		const { treeData } = get();

		set({
			treeData: produce(treeData, (draft) => {
				const node = findNode(draft, nodeId);
				if (node) {
					node.dirinfo = dirinfo;
				}
			}),
		});
	},

	refreshNodeHierarchy: async (nodeId) => {
		const { treeData } = get();
		const node = findNode(treeData, nodeId);
		if (!node) return;

		const path = [];
		let current = node;
		while (current.parent) {
			path.unshift(current.parent.id);
			current = current.parent;
		}

		for (const id of [...path, nodeId]) {
			await get().loadNode(id, 0);
		}
	},

	refreshChildren: async (nodeId) => {
		const { treeData } = get();
		const node = findNode(treeData, nodeId);
		if (!node || !node.children) return;

		const refreshNodeAndChildren = async (currentNode) => {
			await get().loadNode(currentNode.id, 0);

			const updatedNode = findNode(treeData, currentNode.id);
			if (!updatedNode || !updatedNode.children) return;

			for (const child of updatedNode.children) {
				await refreshNodeAndChildren(child);
			}
		};

		for (const child of node.children) {
			await refreshNodeAndChildren(child);
		}
	},

	refreshNodeComplete: async (nodeId) => {
		const { refreshNodeHierarchy, refreshChildren } = get();
		await refreshNodeHierarchy(nodeId);
		await refreshChildren(nodeId);
	},
}));
