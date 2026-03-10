package hu.bme.aut.gallery

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import hu.bme.aut.gallery.databinding.ActivityGalleryBinding

class GalleryActivity : AppCompatActivity() {

    private lateinit var binding: ActivityGalleryBinding

    override fun onCreate(savedInstanceState: Bundle?) {

        try {
            Thread.sleep(1000)
        } catch (e: InterruptedException) {
            e.printStackTrace()
        }
        // hide action bar
        supportActionBar?.hide()

        super.onCreate(savedInstanceState)
        binding = ActivityGalleryBinding.inflate(layoutInflater)
        setContentView(binding.root)


        binding.folderButton.setOnClickListener {
            startActivity(Intent(this, FolderActivity::class.java))
        }
        binding.calendarButton.setOnClickListener {
            startActivity(Intent(this, CalendarActivity::class.java))
        }


    }
}