import android.Manifest.permission
import android.R
import android.content.pm.PackageManager
import android.os.Bundle
import android.os.Environment
import android.provider.MediaStore
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import androidx.recyclerview.widget.GridLayoutManager
import androidx.recyclerview.widget.RecyclerView

class MainActivity : AppCompatActivity() {
    private var imagePaths: ArrayList<String>? = null
    private var imagesRV: RecyclerView? = null
    private var imageRVAdapter: RecyclerViewAdapter? = null
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)


        // creating a new array list and
        // initializing our recycler view.
        imagePaths = ArrayList()
        imagesRV = findViewById(R.id.idRVImages)

        // we are calling a method to request
        // the permissions to read external storage.
        requestPermissions()

        // calling a method to
        // prepare our recycler view.
        prepareRecyclerView()
    }

    private fun checkPermission(): Boolean {
        // in this method we are checking if the permissions are granted or not and returning the result.
        val result =
            ContextCompat.checkSelfPermission(applicationContext, permission.READ_EXTERNAL_STORAGE)
        return result == PackageManager.PERMISSION_GRANTED
    }

    private fun requestPermissions() {
        if (checkPermission()) {
            // if the permissions are already granted we are calling
            // a method to get all images from our external storage.
            Toast.makeText(this, "Permissions granted..", Toast.LENGTH_SHORT).show()
            imagePath
        } else {
            // if the permissions are not granted we are
            // calling a method to request permissions.
            requestPermission()
        }
    }

    private fun requestPermission() {
        //on below line we are requesting the read external storage permissions.
        ActivityCompat.requestPermissions(
            this,
            arrayOf(permission.READ_EXTERNAL_STORAGE),
            PERMISSION_REQUEST_CODE
        )
    }

    private fun prepareRecyclerView() {

        // in this method we are preparing our recycler view.
        // on below line we are initializing our adapter class.
        imageRVAdapter = RecyclerViewAdapter(this@MainActivity, imagePaths)

        // on below line we are creating a new grid layout manager.
        val manager = GridLayoutManager(this@MainActivity, 4)

        // on below line we are setting layout
        // manager and adapter to our recycler view.
        imagesRV!!.layoutManager = manager
        imagesRV!!.adapter = imageRVAdapter
    }// on below line we are moving our cursor position

    // on below line we are getting image file path

    // after that we are getting the image file path
    // and adding that path in our array list.
    // after adding the data to our
    // array list we are closing our cursor.
// if the sd card is present we are creating a new list in
    // which we are getting our images data with their ids.

    // on below line we are creating a new
    // string to order our images by string.

    // this method will stores all the images
    // from the gallery in Cursor

    // below line is to get total number of images

    // on below line we are running a loop to add
    // the image file path in our array list.
    // in this method we are adding all our image paths
    // in our arraylist which we have created.
    // on below line we are checking if the device is having an sd card or not.
    private val imagePath: Unit
        private get() {
            // in this method we are adding all our image paths
            // in our arraylist which we have created.
            // on below line we are checking if the device is having an sd card or not.
            val isSDPresent = Environment.getExternalStorageState() == Environment.MEDIA_MOUNTED
            if (isSDPresent) {

                // if the sd card is present we are creating a new list in
                // which we are getting our images data with their ids.
                val columns = arrayOf(MediaStore.Images.Media.DATA, MediaStore.Images.Media._ID)

                // on below line we are creating a new
                // string to order our images by string.
                val orderBy = MediaStore.Images.Media._ID

                // this method will stores all the images
                // from the gallery in Cursor
                val cursor = contentResolver.query(
                    MediaStore.Images.Media.EXTERNAL_CONTENT_URI,
                    columns,
                    null,
                    null,
                    orderBy
                )

                // below line is to get total number of images
                val count = cursor!!.count

                // on below line we are running a loop to add
                // the image file path in our array list.
                for (i in 0 until count) {

                    // on below line we are moving our cursor position
                    cursor.moveToPosition(i)

                    // on below line we are getting image file path
                    val dataColumnIndex = cursor.getColumnIndex(MediaStore.Images.Media.DATA)

                    // after that we are getting the image file path
                    // and adding that path in our array list.
                    imagePaths!!.add(cursor.getString(dataColumnIndex))
                }
                imageRVAdapter.notifyDataSetChanged()
                // after adding the data to our
                // array list we are closing our cursor.
                cursor.close()
            }
        }

    override fun onRequestPermissionsResult(
        requestCode: Int,
        permissions: Array<String>,
        grantResults: IntArray
    ) {
        // this method is called after permissions has been granted.
        when (requestCode) {
            PERMISSION_REQUEST_CODE ->                // in this case we are checking if the permissions are accepted or not.
                if (grantResults.size > 0) {
                    val storageAccepted = grantResults[0] == PackageManager.PERMISSION_GRANTED
                    if (storageAccepted) {
                        // if the permissions are accepted we are displaying a toast message
                        // and calling a method to get image path.
                        Toast.makeText(this, "Permissions Granted..", Toast.LENGTH_SHORT).show()
                        imagePath
                    } else {
                        // if permissions are denied we are closing the app and displaying the toast message.
                        Toast.makeText(
                            this,
                            "Permissions denied, Permissions are required to use the app..",
                            Toast.LENGTH_SHORT
                        ).show()
                    }
                }
        }
    }

    companion object {
        // on below line we are creating variables for
        // our array list, recycler view and adapter class.
        private const val PERMISSION_REQUEST_CODE = 200
    }
}