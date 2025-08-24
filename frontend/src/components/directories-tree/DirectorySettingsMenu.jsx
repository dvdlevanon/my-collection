import { Menu, MenuItem, Typography } from '@mui/material';
import Client from '../../utils/client';
import PathUtil from '../../utils/path-util';

function DirectorySettingsMenu({ anchorEl, anchorPosition, refreshDir, dir, onClose }) {
	const includeDir = () => {
		Client.includeDir(dir.id, false, false).then(() => {
			refreshDir(dir);
			onClose();
		});
	};

	const includeSubdirs = () => {
		Client.includeDir(dir.id, true, false).then(() => {
			refreshDir(dir);
			onClose();
		});
	};

	const includeHierarchy = () => {
		Client.includeDir(dir.id, false, true).then(() => {
			refreshDir(dir);
			onClose();
		});
	};

	const excludeDir = () => {
		Client.excludeDir(dir.id).then(() => {
			refreshDir(dir);
			onClose();
		});
	};

	const included = () => {
		return dir.dirinfo && !dir.dirinfo.excluded;
	};

	return (
		<Menu
			open={anchorEl != null}
			anchorEl={anchorEl}
			anchorPosition={anchorPosition}
			anchorReference="anchorPosition"
			onClose={onClose}
		>
			<MenuItem disabled>
				<Typography variant="h5" color="white">
					{'Actions for:  "' + PathUtil.dirname(dir.name) + '"'}
				</Typography>
			</MenuItem>
			{!included() && <MenuItem onClick={includeDir}>Include</MenuItem>}
			{!included() && <MenuItem onClick={includeSubdirs}>Include With Sub Directories</MenuItem>}
			{!included() && <MenuItem onClick={includeHierarchy}>Include With All Hierarchy</MenuItem>}
			{included() && <MenuItem onClick={excludeDir}>Excluded Hierarchy</MenuItem>}
		</Menu>
	);
}

export default DirectorySettingsMenu;
