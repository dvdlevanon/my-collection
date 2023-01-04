function Tag({tag, onTagActivated}) {
  return (
    <>
        <div className={tag.selected ? "tag selected" : "tag"} onClick={() => onTagActivated(tag)}>
            {tag.title}
        </div>
    </>
  )
}

export default Tag