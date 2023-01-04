function GalleryItem({ item }) {
  return (
    <>
        <img src={item.cover} alt='' />
        <span>{item.title}</span>
    </>
  )
}

export default GalleryItem