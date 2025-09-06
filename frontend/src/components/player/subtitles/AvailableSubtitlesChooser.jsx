import { FormControl, MenuItem, Select, Stack } from '@mui/material';
import { useQuery } from '@tanstack/react-query';
import { useEffect, useRef, useState } from 'react';
import ReactQueryUtil from '../../../utils/react-query-util';
import { usePlayerStore } from '../PlayerStore';
import { useSubtitleStore } from './SubtitlesStore';

const NO_SUBTITLE = 'No Subtitles Available';

function AvailableSubtitlesChooser() {
	const containerRef = useRef();
	const playerStore = usePlayerStore();
	const subtitleStore = useSubtitleStore();
	const availableSubtitleQuery = useQuery(ReactQueryUtil.availableSubtitleQuery(playerStore.itemId));
	const [selectedSubtitle, setSelectedSubtitle] = useState(NO_SUBTITLE);

	useEffect(() => {
		if (availableSubtitleQuery.data) {
		}
	}, [availableSubtitleQuery.data]);

	const getSubtitles = () => {
		if (!availableSubtitleQuery.data) {
			return [NO_SUBTITLE];
		}

		if (selectedSubtitle == NO_SUBTITLE) {
			setSelectedSubtitle(availableSubtitleQuery.data[0]);
		}

		return availableSubtitleQuery.data;
	};

	const onChange = (event) => {
		let value = event.target.value;

		setSelectedSubtitle(value);
		subtitleStore.setSelectedSubtitleName(value);
	};

	return (
		<Stack ref={containerRef}>
			<FormControl fullWidth>
				<Select
					onChange={onChange}
					value={selectedSubtitle}
					MenuProps={{
						container: containerRef.current,
						PaperProps: {
							container: containerRef.current,
						},
					}}
				>
					{getSubtitles().map((subtitle) => {
						return (
							<MenuItem key={subtitle} value={subtitle}>
								{subtitle}
							</MenuItem>
						);
					})}
				</Select>
			</FormControl>
		</Stack>
	);
}

export default AvailableSubtitlesChooser;
